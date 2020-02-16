/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	mutablehome "github.com/djthorpe/mutablehome"
	dvb "github.com/djthorpe/mutablehome/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Demux struct {
	Adapter  uint
	Demux    uint
	Frontend mutablehome.DVBFrontend
	FilePoll gopi.FilePoll
	Bus      gopi.Bus
}

type demux struct {
	adapter, demux uint
	frontend       mutablehome.DVBFrontend
	filepoll       gopi.FilePoll
	bus            gopi.Bus
	sectionfilter  map[uintptr]*SectionFilter
	streamfilter   *StreamFilter

	base.Unit
	sync.RWMutex // Used for access to filters
	sync.Mutex   // Used for method access
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Demux) Name() string { return "mutablehome/dvb/demux" }

func (config Demux) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(demux)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.DVBDemux

func (this *demux) Init(config Demux) error {
	// Check for demux device
	path := dvb.DVB_DMXPath(config.Adapter, config.Demux)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return gopi.ErrBadParameter.WithPrefix("dvb.demux")
	} else {
		this.demux = config.Demux
	}

	// Check for frontend device
	if config.Frontend == nil {
		return gopi.ErrBadParameter.WithPrefix("dvb.adapter")
	} else {
		this.frontend = config.Frontend
	}

	// Check for filepoll
	if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("filepoll")
	} else {
		this.filepoll = config.FilePoll
	}

	// Check for bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}

	// Create filter map
	this.sectionfilter = make(map[uintptr]*SectionFilter)

	// Return success
	return nil
}

func (this *demux) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	errs := gopi.NewCompoundError()

	if this.streamfilter != nil {
		errs.Add(this.filepoll.Unwatch(this.streamfilter.Fd()))
		errs.Add(this.streamfilter.Close())
		this.setStreamFilter(nil)
	}

	for fd, filter := range this.sectionfilter {
		errs.Add(this.filepoll.Unwatch(fd))
		errs.Add(filter.Close())
		this.setSectionFilterForFd(fd, nil)
	}

	// Release resources
	this.sectionfilter = nil
	this.streamfilter = nil
	this.frontend = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *demux) NewSectionFilter(pid uint16, tid mutablehome.DVBTableType) (mutablehome.DVBFilter, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Create filter
	filter, err := NewSectionFilter(this.adapter, this.demux, pid, tid)
	if err != nil {
		return nil, err
	}

	// Watch filter
	if err := this.filepoll.Watch(filter.Fd(), gopi.FILEPOLL_FLAG_READ, this.Read); err != nil {
		filter.Close()
		return nil, err
	}

	// Set filter
	this.setSectionFilterForFd(filter.Fd(), filter)

	// Start filtering
	if err := filter.Start(); err != nil {
		this.filepoll.Unwatch(filter.Fd())
		filter.Close()
		this.setSectionFilterForFd(filter.Fd(), nil)
		return nil, err
	}

	// Return success
	return filter, nil
}

func (this *demux) ScanPAT() (mutablehome.DVBFilter, error) {
	return this.NewSectionFilter(uint16(0x00), mutablehome.DVB_TS_TABLE_PAT)
}

func (this *demux) ScanSDT(other bool) (mutablehome.DVBFilter, error) {
	if other {
		return this.NewSectionFilter(uint16(0x11), mutablehome.DVB_TS_TABLE_SDT_OTHER)
	} else {
		return this.NewSectionFilter(uint16(0x11), mutablehome.DVB_TS_TABLE_SDT)
	}
}

func (this *demux) ScanNIT(other bool) (mutablehome.DVBFilter, error) {
	if other {
		return this.NewSectionFilter(uint16(0x10), mutablehome.DVB_TS_TABLE_NIT_OTHER)
	} else {
		return this.NewSectionFilter(uint16(0x10), mutablehome.DVB_TS_TABLE_NIT)
	}
}

func (this *demux) ScanEITNowNext(other bool) (mutablehome.DVBFilter, error) {
	if other {
		return this.NewSectionFilter(uint16(0x12), mutablehome.DVB_TS_TABLE_EIT_OTHER)
	} else {
		return this.NewSectionFilter(uint16(0x12), mutablehome.DVB_TS_TABLE_EIT)
	}
}

func (this *demux) ScanPMT(section mutablehome.DVBSection) ([]mutablehome.DVBFilter, error) {
	if pat, ok := section.(*SectionPAT); ok == false || pat == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("section")
	} else {
		filters := make([]mutablehome.DVBFilter, 0, len(pat.Programs))
		for _, row := range pat.Programs {
			if row.Program == 0 {
				// Ignore NIT
				continue
			} else if filter, err := this.NewSectionFilter(row.Pid, mutablehome.DVB_TS_TABLE_PMT); err != nil {
				return nil, err
			} else {
				filters = append(filters, filter)
			}
		}
		return filters, nil
	}
}

func (this *demux) DestroyFilter(filter mutablehome.DVBFilter) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	errs := gopi.NewCompoundError()

	if fd := filter.Fd(); fd == 0 {
		return gopi.ErrBadParameter.WithPrefix("filter")
	} else if filter_, ok := filter.(*StreamFilter); ok && filter_ == this.streamfilter {
		errs.Add(this.filepoll.Unwatch(fd))
		errs.Add(filter_.Close())
		this.setStreamFilter(nil)
	} else if filter_ := this.sectionFilterForFd(fd); filter_ != nil {
		errs.Add(this.filepoll.Unwatch(fd))
		errs.Add(filter_.Close())
		this.setSectionFilterForFd(fd, nil)
	} else {
		return gopi.ErrBadParameter.WithPrefix("filter")
	}

	// Return any error
	return errs.ErrorOrSelf()
}

func (this *demux) NewStreamFilter(pids []uint16) (mutablehome.DVBFilter, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.streamfilter != nil {
		return nil, gopi.ErrOutOfOrder.WithPrefix("filter")
	} else if len(pids) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("pids")
	} else if filter, err := NewStreamFilter(this.adapter, this.demux, pids[0], dvb.DVB_DMX_IN_FRONTEND, dvb.DVB_DMX_OUT_TSDEMUX_TAP, dvb.DVB_DMX_PES_OTHER); err != nil {
		return nil, err
	} else if err := filter.AddPids(pids[1:]); err != nil {
		filter.Close()
		return nil, err
	} else if err := this.filepoll.Watch(filter.Fd(), gopi.FILEPOLL_FLAG_READ, this.Read); err != nil {
		filter.Close()
		return nil, err
	} else {
		this.setStreamFilter(filter)
		this.streamfilter = filter
		if err := filter.Start(); err != nil {
			filter.Close()
			this.streamfilter = nil
			return nil, err
		} else {
			return filter, nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *demux) sectionFilterForFd(fd uintptr) *SectionFilter {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if filter, exists := this.sectionfilter[fd]; exists {
		return filter
	} else {
		return nil
	}
}

func (this *demux) setSectionFilterForFd(fd uintptr, filter *SectionFilter) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if filter != nil {
		this.sectionfilter[fd] = filter
	} else {
		delete(this.sectionfilter, fd)
	}
}

func (this *demux) setStreamFilter(filter *StreamFilter) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.streamfilter = filter
}

func (this *demux) streamFilterForFd(fd uintptr) *StreamFilter {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.streamfilter == nil {
		return nil
	} else if this.streamfilter.Fd() == fd {
		return this.streamfilter
	} else {
		return nil
	}
}

func (this *demux) Read(fd uintptr, flags gopi.FilePollFlags) {
	if flags&gopi.FILEPOLL_FLAG_READ == gopi.FILEPOLL_FLAG_READ {
		if filter := this.streamFilterForFd(fd); filter != nil {
			if buf, err := NewTSReaderRead(fd, TS_PACKET_LENGTH*4); err != nil {
				this.Log.Warn("Stream Read error:", err)
			} else {
				this.Log.Debug("buf=", buf.Size(), "bytes")
			}
			return
		} else if filter := this.sectionFilterForFd(fd); filter != nil {
			if section, err := TSRead(fd); err != nil {
				this.Log.Warn("Section Read error:", err)
				return
			} else {
				this.bus.Emit(NewSectionEvent(this, filter, section))
				return
			}
		}
	}
	this.Log.Warn("Invalid file descriptor:", fd)
}
