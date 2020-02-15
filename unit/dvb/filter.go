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
	home "github.com/djthorpe/mutablehome"
	dvb "github.com/djthorpe/mutablehome/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Filter struct {
	dev *os.File
	sync.Mutex
}

type SectionFilter struct {
	Filter
	dvb.DMXSectionFilter
}

type StreamFilter struct {
	Filter
	dvb.DMXStreamFilter
}

////////////////////////////////////////////////////////////////////////////////
// NEW / CLOSE

func NewSectionFilter(adapter, demux uint, pid uint16, tid home.DVBTableType) (*SectionFilter, error) {
	if dev, err := dvb.DVB_DMXOpen(adapter, demux); err != nil {
		return nil, err
	} else {
		filter := &SectionFilter{
			Filter{
				dev: dev,
			},
			dvb.DMXSectionFilter{
				Pid:     pid,
				Timeout: 0,
				Flags:   dvb.DVB_DMX_FLAG_IMMEDIATE_START,
			},
		}
		filter.DMXSectionFilter.Pattern.Filter[0] = uint8(tid)
		filter.DMXSectionFilter.Pattern.Mask[0] = 0xFF

		if err := dvb.DVB_DMXSetSectionFilter(dev.Fd(), filter.DMXSectionFilter); err != nil {
			dev.Close()
			return nil, err
		} else {
			return filter, nil
		}

	}
}

func NewStreamFilter(adapter, demux uint, pid uint16, in dvb.DMXInput, out dvb.DMXOutput, streamType dvb.DMXStreamType) (*StreamFilter, error) {
	if dev, err := dvb.DVB_DMXOpen(adapter, demux); err != nil {
		return nil, err
	} else {
		filter := &StreamFilter{
			Filter{
				dev: dev,
			},
			dvb.DMXStreamFilter{
				Pid:  pid,
				In:   in,
				Out:  out,
				Type: streamType,
			},
		}
		if err := dvb.DVB_DMXSetStreamFilter(dev.Fd(), filter.DMXStreamFilter); err != nil {
			dev.Close()
			return nil, err
		} else {
			return filter, nil
		}
	}
}

func (this *Filter) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev != nil {
		return this.dev.Close()
	}

	// Release resources
	this.dev = nil

	// Return success
	return nil
}

func (this *Filter) Fd() uintptr {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return 0
	} else {
		return this.dev.Fd()
	}
}

func (this *Filter) Start() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		return dvb.DVB_DMXStart(this.dev.Fd())
	}
}

func (this *Filter) Stop() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		return dvb.DVB_DMXStop(this.dev.Fd())
	}
}

func (this *Filter) AddPid(pid uint16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		return dvb.DVB_DMXAddPid(this.dev.Fd(), pid)
	}
}

func (this *Filter) AddPids(pids []uint16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		for _, pid := range pids {
			if err := dvb.DVB_DMXAddPid(this.dev.Fd(), pid); err != nil {
				return err
			}
		}
	}
	// Success
	return nil
}

func (this *Filter) SetBufferSize(size uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		return dvb.DVB_DMXSetBufferSize(this.dev.Fd(), size)
	}
}

func (this *Filter) RemovePid(pid uint16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else {
		return dvb.DVB_DMXRemovePid(this.dev.Fd(), pid)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func openStreamFilter(adapter, demux uint, out dvb.DMXOutput) (*os.File, error) {
	switch out {
	case dvb.DVB_DMX_OUT_TS_TAP:
		if dev, err := dvb.DVB_DVROpen(adapter, demux); err != nil {
			return nil, err
		} else {
			return dev, nil
		}
	case dvb.DVB_DMX_OUT_TAP, dvb.DVB_DMX_OUT_TSDEMUX_TAP:
		if dev, err := dvb.DVB_DMXOpen(adapter, demux); err != nil {
			return nil, err
		} else {
			return dev, nil
		}
	default:
		return nil, gopi.ErrBadParameter.WithPrefix("out")
	}
}
