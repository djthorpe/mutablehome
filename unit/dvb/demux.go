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
}

type demux struct {
	adapter, demux uint
	frontend       mutablehome.DVBFrontend
	filter         map[uintptr]filter

	base.Unit
	sync.Mutex
}

type sectionfilter struct {
	filter
}

type streamfilter struct {
	filter
}

type filter struct {
	dev *os.File

	sync.Mutex
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

	// Create filter map
	this.filter = make(map[uintptr]filter)

	// Return success
	return nil
}

func (this *demux) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	errs := gopi.NewCompoundError()
	for fd, filter := range this.filter {
		errs.Add(filter.Close())
		delete(this.filter, fd)
	}

	// Release resources
	this.filter = nil
	this.frontend = nil

	// Return success
	return this.Unit.Close()
}

func (this *demux) NewSectionFilter() (mutablehome.DVBFilter, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if dev, err := dvb.DVB_DMXOpen(this.adapter, this.demux); err != nil {
		return nil, err
	} else {
		this.filter[dev.Fd()] = filter{dev: dev}
		return this.filter[dev.Fd()], nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION filter

func (this *filter) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev != nil {
		return this.dev.Close()
	} else {
		return nil
	}
}
