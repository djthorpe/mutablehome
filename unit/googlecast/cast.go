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
}

type cast struct {
	log       gopi.Logger
	discovery gopi.RPCServiceDiscovery
	bus       gopi.Bus

	devices
	lookup
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
		this.lookup.Discovery = config.Discovery
	}

	// Check for bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}

	// Init
	this.devices.Init()

	// Handle bus messages
	if err := this.bus.NewHandler(gopi.EventHandler{Name: "gopi.RPCEvent", Handler: this.EventHandler}); err != nil {
		return err
	}

	// Start discovery
	if err := this.lookup.Start(SERVICE_TYPE_GOOGLECAST); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *cast) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Stop discovery
	this.lookup.Stop()

	// Release device resources
	this.devices.Close()

	// Release resources
	this.bus = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Devices

func (this *cast) Devices() []iface.CastDevice {
	devices := make([]iface.CastDevice, 0, len(this.devices.devices))
	for _, device := range this.devices.devices {
		devices = append(devices, device)
	}
	return devices
}

func (this *cast) Connect(iface.CastDevice, gopi.RPCFlag) error {
	return gopi.ErrNotImplemented
}

func (this *cast) Disconnect(iface.CastDevice) error {
	return gopi.ErrNotImplemented
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
		if device, updated := this.devices.Update(evt.Service()); updated == true {
			this.bus.Emit(NewAddedEvent(this, device))
		}
	case gopi.RPC_EVENT_SERVICE_UPDATED:
		if device, updated := this.devices.Update(evt.Service()); updated == true {
			this.bus.Emit(NewUpdatedEvent(this, device))
		}
	case gopi.RPC_EVENT_SERVICE_EXPIRED, gopi.RPC_EVENT_SERVICE_REMOVED:
		if device, removed := this.devices.Remove(evt.Service()); removed == true {
			this.bus.Emit(NewRemovedEvent(this, device))
		}
	default:
		this.Log.Warn("Unhandled event", evt.Type(), "for", evt.Service().Name)
	}
}
