/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	mutablehome "github.com/djthorpe/mutablehome"
	dvb "github.com/djthorpe/mutablehome/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Frontend struct {
	Adapter  uint
	Frontend uint
}

type frontend struct {
	dev     *os.File
	version string
	name    string
	systems []mutablehome.DVBDeliverySystem

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Frontend) Name() string { return "mutablehome/dvb/frontend" }

func (config Frontend) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(frontend)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.DVBFrontend

func (this *frontend) Init(config Frontend) error {
	if dev, err := dvb.DVB_FEOpen(config.Adapter, config.Frontend); err != nil {
		return err
	} else {
		this.dev = dev
	}
	if major, minor, err := dvb.DVB_FEVersion(this.dev.Fd()); err != nil {
		return err
	} else {
		this.version = fmt.Sprintf("%v.%v", major, minor)
	}
	if info, err := dvb.DVB_FEGetInfo(this.dev.Fd()); err != nil {
		return err
	} else {
		this.name = info.Name()
	}
	if systems, err := dvb.DVB_FEDeliverySystemEnum(this.dev.Fd()); err != nil {
		return err
	} else {
		this.systems = systems
	}

	// Return success
	return nil
}

func (this *frontend) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.dev = nil
	this.systems = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *frontend) Name() string {
	return this.name
}

func (this *frontend) DeliverySystems() []mutablehome.DVBDeliverySystem {
	return this.systems
}

////////////////////////////////////////////////////////////////////////////////
// TUNE

func (this *frontend) Tune(mutablehome.DVBProperties) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *frontend) String() string {
	if this.dev == nil {
		return "<" + this.Log.Name() + ">"
	} else {
		return "<" + this.Log.Name() +
			" path=" + strconv.Quote(this.dev.Name()) +
			" name=" + strconv.Quote(this.Name()) +
			" api_version=" + strconv.Quote(this.version) +
			" delivery_systems=" + fmt.Sprint(this.systems) +
			">"
	}
}
