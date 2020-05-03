/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package node

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type node struct {
	base.Unit
	base.PubSub
	sync.Mutex
	sync.WaitGroup

	gateway mutablehome.TradfriGateway
	stop    chan struct{}
	devices map[uint]*device
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Update devices, groups and scenes regularly
	DELTA_INTERVAL = 15 * time.Second
)

/*

type nodedevice struct {
	gateway mutablehome.TradfriGateway
	device  mutablehome.TradfriDevice
}
*/

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *node) Init(config Node) error {
	// Set up gateway
	if config.Gateway == nil {
		return gopi.ErrBadParameter.WithPrefix("tradfri")
	} else {
		this.gateway = config.Gateway
	}

	// Create stop signal
	this.stop = make(chan struct{})
	this.devices = make(map[uint]*device)

	// Success
	return nil
}

func (this *node) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close stop channel and wait for background process to end
	close(this.stop)
	this.WaitGroup.Wait()

	// Release resources
	this.devices = nil
	this.gateway = nil
	this.stop = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *node) String() string {
	str := "<" + this.Log.Name()
	if id := this.Id(); id != "" {
		str += " id=" + strconv.Quote(id)
	}
	str += " name=" + strconv.Quote(this.Name())
	str += " gateway=" + fmt.Sprint(this.gateway)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT AND DISCONNECT GATEWAY

func (this *node) Connect(service gopi.RPCServiceRecord, flags gopi.RPCFlag) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := this.gateway.Connect(service, flags); err != nil {
		return err
	} else {
		go this.BackgroundProcess(this.stop)
	}
	this.Emit(this.NewGatewayEvent(mutablehome.EVENT_NODE_ONLINE))
	return nil
}

func (this *node) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := this.gateway.Disconnect(); err != nil {
		return err
	}
	this.Emit(this.NewGatewayEvent(mutablehome.EVENT_NODE_OFFLINE))
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *node) Id() string {
	return this.gateway.Id()
}

func (this *node) Name() string {
	v := this.gateway.Version()
	if v != "" {
		return "Ikea Tradfri Gateway " + v
	} else {
		return "Ikea Tradfri Gateway (offline)"
	}
}

func (this *node) Device(key string) mutablehome.Device {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if id, err := strconv.ParseUint(key, 10, 64); err != nil {
		return nil
	} else if device, exists := this.devices[uint(id)]; exists == false {
		return nil
	} else {
		return device
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESS

func (this *node) BackgroundProcess(stop <-chan struct{}) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	this.Log.Debug("Start of background process")
	ticker := time.NewTimer(500 * time.Millisecond)
FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			if err := this.BackgroundDiscoverDevices(); err != nil {
				this.Log.Error(err)
			}
			ticker.Reset(DELTA_INTERVAL)
		case <-stop:
			ticker.Stop()
			break FOR_LOOP
		}
	}
	this.Log.Debug("End of background process")
}

func (this *node) BackgroundDiscoverDevices() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if devices, err := this.gateway.Devices(); err != nil {
		return fmt.Errorf("DiscoverDevices: %w", err)
	} else {
		for _, device := range this.devices {
			device.seen = false
		}
		for _, id := range devices {
			if device, exists := this.devices[id]; exists {
				device.seen = true
			} else {
				fmt.Println("TODO: New Device", id)
			}
		}
		for _, device := range this.devices {
			if device.seen == false {
				fmt.Println("TODO: Remove Device", device.id)
			}
		}
	}

	// Return success
	return nil
}

/*

func (this *nodedevice) Power() mutablehome.TraitType {
	// If not a light, then return unknown
	if this.device.Type() != mutablehome.IKEA_DEVICE_TYPE_LIGHT {
		return mutablehome.TRAIT_NONE
	}
	// If no lights, then return unknown
	if len(this.device.Lights()) == 0 {
		return mutablehome.TRAIT_NONE
	}
	// Take first light value
	if light := this.device.Lights()[0]; light.Power() == true {
		return mutablehome.TRAIT_POWER_ON
	} else {
		return mutablehome.TRAIT_POWER_OFF
	}
}

func (this *nodedevice) SetPower(state mutablehome.TraitType) error {
	// If not a light, then return error
	if this.device.Type() != mutablehome.IKEA_DEVICE_TYPE_LIGHT {
		return gopi.ErrNotImplemented.WithPrefix("Power")
	}
	// If no lights, then return error
	if len(this.device.Lights()) == 0 {
		return gopi.ErrNotImplemented.WithPrefix("Power")
	}
	// Set power value for all lights
	for _, light := range this.device.Lights() {
		switch state {
		case mutablehome.TRAIT_POWER_ON:
			return this.gateway.Send(light.SetPower(true))
		case mutablehome.TRAIT_POWER_OFF:
			return this.gateway.Send(light.SetPower(false))
		default:
			return gopi.ErrBadParameter.WithPrefix("Power")
		}
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// EVENT HANDLER

func (this *node) EventHandlerFunc(_ context.Context, _ gopi.App, evt gopi.Event) {
	evt_ := evt.(mutablehome.IkeaEvent)
	device_ := &nodedevice{evt.Source().(mutablehome.Ikea), evt_.Device()}
	switch evt_.Type() {
	case mutablehome.IKEA_EVENT_GATEWAY_CONNECTED:
		this.Emit(this.NewEvent(mutablehome.EVENT_NODE_ONLINE, nil))
	case mutablehome.IKEA_EVENT_GATEWAY_DISCONNECTED:
		this.Emit(this.NewEvent(mutablehome.EVENT_NODE_OFFLINE, nil))
	case mutablehome.IKEA_EVENT_DEVICE_ADDED:
		this.Emit(this.NewEvent(mutablehome.EVENT_DEVICE_ADDED, device_))
	case mutablehome.IKEA_EVENT_DEVICE_CHANGED:
		this.Emit(this.NewEvent(mutablehome.EVENT_DEVICE_METADATA_CHANGED, device_))
	default:
		this.Log.Warn("Ignoring:", evt_)
	}
}
*/
