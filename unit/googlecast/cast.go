/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"context"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	iface "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cast struct {
	Discovery gopi.RPCServiceDiscovery
	Bus       gopi.Bus
	Timeout   time.Duration
}

type cast struct {
	log       gopi.Logger
	discovery gopi.RPCServiceDiscovery
	bus       gopi.Bus
	timeout   time.Duration
	devices   map[string]*device

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// COMSTANTS

const (
	SERVICE_TYPE_GOOGLECAST = "_googlecast._tcp"
	DELTA_LOOKUP_TIME       = 60 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Cast) Name() string { return "googlecast" }

func (config Cast) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(cast)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *cast) Init(config Cast) error {
	// Check for discovery
	if config.Discovery == nil {
		return gopi.ErrBadParameter.WithPrefix("discovery")
	} else {
		this.devices = make(map[string]*device)
		this.discovery = config.Discovery
		this.timeout = config.Timeout
	}

	// Check for bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}

	// Handle bus messages
	if err := this.bus.NewHandler(gopi.EventHandler{Name: "gopi.RPCEvent", Handler: this.EventHandler}); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *cast) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Disconnect devices
	for _, device := range this.devices {
		if err := device.Disconnect(); err != nil {
			this.Log.Warn(err)
		}
	}

	// Release resources
	this.devices = nil
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Devices

func (this *cast) Devices(ctx context.Context) ([]iface.CastDevice, error) {

	// Perform the lookup
	if _, err := this.discovery.Lookup(ctx, SERVICE_TYPE_GOOGLECAST); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		return nil, err
	}

	// Lock for reading devices
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Enumerate the devices
	devices := make([]iface.CastDevice, 0, len(this.devices))
	for _, device := range this.devices {
		devices = append(devices, device)
	}

	// Return success
	return devices, nil
}

func (this *cast) Connect(d iface.CastDevice, flags gopi.RPCFlag) error {
	// Check parameters
	if d == nil {
		return gopi.ErrBadParameter.WithPrefix("device")
	}

	// Typecast and make connection
	if d_, ok := d.(*device); ok == false {
		return gopi.ErrBadParameter.WithPrefix("device")
	} else if err := d_.Connect(flags, this.timeout); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *cast) Disconnect(d iface.CastDevice) error {
	// Check parameters
	if d == nil {
		return gopi.ErrBadParameter.WithPrefix("device")
	}

	// Typecast and do disconnection
	if d_, ok := d.(*device); ok == false {
		return gopi.ErrBadParameter.WithPrefix("device")
	} else {
		return d_.Disconnect()
	}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT HANDLER

func (this *cast) EventHandler(_ context.Context, _ gopi.App, evt gopi.Event) {
	if evt_ := evt.(gopi.RPCEvent); evt_.Service().Service == SERVICE_TYPE_GOOGLECAST {
		this.RPCEventHandler(evt_)
	}
}

func (this *cast) RPCEventHandler(evt gopi.RPCEvent) {
	switch evt.Type() {
	case gopi.RPC_EVENT_SERVICE_ADDED:
		if device, updated := this.UpdateDevice(evt.Service()); updated == true {
			this.bus.Emit(NewAddedEvent(this, device))
		}
	case gopi.RPC_EVENT_SERVICE_UPDATED:
		if evt.Service().Port > 0 && evt.Service().Host != "" {
			if device, updated := this.UpdateDevice(evt.Service()); updated == true {
				this.bus.Emit(NewUpdatedEvent(this, device))
			}
		}
	case gopi.RPC_EVENT_SERVICE_EXPIRED, gopi.RPC_EVENT_SERVICE_REMOVED:
		if device, removed := this.RemoveDevice(evt.Service()); removed == true {
			this.bus.Emit(NewRemovedEvent(this, device))
			if err := device.Disconnect(); err != nil {
				this.Log.Error(err)
			}
		}
	default:
		this.Log.Warn("Unhandled event", evt.Type(), "for", evt.Service().Name)
	}
}

////////////////////////////////////////////////////////////////////////////////
// ADD, REMOVE AND UPDATE DEVICES

func (this *cast) UpdateDevice(srv gopi.RPCServiceRecord) (*device, bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.devices == nil {
		return nil, false
	} else if key := srv.Name; key == "" {
		return nil, false
	} else if d, exists := this.devices[key]; exists {
		d.setService(srv)
		return d, true
	} else if d, err := gopi.New(Device{srv}, this.Log.Clone(Device{}.Name())); err != nil {
		this.Log.Error(err)
		return nil, false
	} else {
		this.devices[key] = d.(*device)
		return d.(*device), true
	}
}

func (this *cast) RemoveDevice(srv gopi.RPCServiceRecord) (*device, bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.devices == nil {
		return nil, false
	} else if key := srv.Name; key == "" {
		return nil, false
	} else if device, exists := this.devices[key]; exists == false {
		return nil, false
	} else {
		delete(this.devices, key)
		return device, true
	}
}
