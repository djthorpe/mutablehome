/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

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

func (this *frontend) Supports(sys mutablehome.DVBDeliverySystem) bool {
	for _, value := range this.systems {
		if sys == value {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// TUNE

func (this *frontend) Tune(ctx context.Context, properties mutablehome.DVBProperties) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Set frontend properties
	if this.dev == nil {
		return gopi.ErrInternalAppError
	} else if properties == nil {
		return gopi.ErrBadParameter.WithPrefix("properties")
	} else if sys, err := properties.DeliverySystem(); err != nil {
		return err
	} else if this.Supports(sys) == false {
		return gopi.ErrNotImplemented.WithPrefix(fmt.Sprint(sys))
	} else if err := dvb.DVB_FEClear(this.dev.Fd()); err != nil {
		return err
	} else {
		switch sys {
		case mutablehome.DVB_SYS_DVBT:
			if err := this.SetPropertiesDVBT(properties); err != nil {
				return err
			}
		default:
			return gopi.ErrNotImplemented.WithPrefix(fmt.Sprint(sys))
		}
	}

	// Begin the tuning
	if err := dvb.DVB_FETune(this.dev.Fd()); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second)
FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			if status, err := dvb.DVB_FEReadStatus(this.dev.Fd()); err != nil {
				return err
			} else {
				fmt.Println(status)
				dvb.DVB_FEStatSignalStrength(this.dev.Fd())
				dvb.DVB_FEStatCarrierNoiseRatio(this.dev.Fd())

				if status&dvb.DVB_FE_STATUS_HAS_LOCK == dvb.DVB_FE_STATUS_HAS_LOCK {
					ticker.Stop()
					break FOR_LOOP
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			break FOR_LOOP
		}
	}

	// Return success
	return ctx.Err()
}

func (this *frontend) SetPropertiesDVBT(properties mutablehome.DVBProperties) error {
	// Set delivery system
	if sys, err := properties.DeliverySystem(); err != nil {
		return err
	} else if err := dvb.DVB_FESetDeliverySystem(this.dev.Fd(), sys); err != nil {
		return err
	}
	// Set frequency
	if freq := properties.Frequency(); freq == 0 {
		return gopi.ErrBadParameter.WithPrefix("frequency")
	} else if err := dvb.DVB_FESetFrequency(this.dev.Fd(), uint(freq)); err != nil {
		return err
	}
	// Set modulation
	if modulation, err := properties.Modulation(); err != nil {
		return err
	} else if err := dvb.DVB_FESetModulation(this.dev.Fd(), modulation); err != nil {
		return err
	}
	// Set bandwidth
	if bandwidth := properties.Bandwidth(); bandwidth == 0 {
		return gopi.ErrBadParameter.WithPrefix("bandwidth")
	} else if err := dvb.DVB_FESetBandwidth(this.dev.Fd(), uint(bandwidth)); err != nil {
		return err
	}
	// Set inversion
	if inversion, err := properties.Inversion(); err != nil {
		return err
	} else if err := dvb.DVB_FESetInversion(this.dev.Fd(), inversion); err != nil {
		return err
	}
	// Set code rate HP
	if codeRate, err := properties.CodeRateHP(); err != nil {
		return err
	} else if err := dvb.DVB_FESetCodeRateHP(this.dev.Fd(), codeRate); err != nil {
		return err
	}
	// Set code rate LP
	if codeRate, err := properties.CodeRateLP(); err != nil {
		return err
	} else if err := dvb.DVB_FESetCodeRateLP(this.dev.Fd(), codeRate); err != nil {
		return err
	}
	// Set guard interval
	if guardInterval, err := properties.GuardInterval(); err != nil {
		return err
	} else if err := dvb.DVB_FESetGuardInterval(this.dev.Fd(), guardInterval); err != nil {
		return err
	}
	// Set transmission mode
	if transmitMode, err := properties.TransmitMode(); err != nil {
		return err
	} else if err := dvb.DVB_FESetTransmitMode(this.dev.Fd(), transmitMode); err != nil {
		return err
	}
	// Set hierarchy
	if hierarchy, err := properties.Hierarchy(); err != nil {
		return err
	} else if err := dvb.DVB_FESetHierarchy(this.dev.Fd(), hierarchy); err != nil {
		return err
	}

	// Return success
	return nil
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
