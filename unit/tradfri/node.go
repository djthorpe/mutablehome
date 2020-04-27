/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (

	// Frameworks
	"context"
	"fmt"

	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Node struct {
	Tradfri     mutablehome.Ikea
	Bus         gopi.Bus
	NodeService NodeService
}

type node struct {
	base.Unit
	base.PubSub

	tradfri mutablehome.Ikea
	bus     gopi.Bus
}

type nodeevent struct {
	Type_   mutablehome.EventType
	Source_ mutablehome.Node
	Device_ mutablehome.Device
}

type nodedevice struct {
	mutablehome.IkeaDevice
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Node) Name() string { return "mutablehome/node/tradfri" }

func (config Node) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(node)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *node) Init(config Node) error {
	if config.Tradfri == nil {
		return gopi.ErrBadParameter.WithPrefix("tradfri")
	} else {
		this.tradfri = config.Tradfri
	}
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}
	if config.NodeService == nil {
		return gopi.ErrBadParameter.WithPrefix("nodeservice")
	} else if err := config.NodeService.SetNode(this); err != nil {
		return err
	}

	// Receive messages of type ikea.Event
	this.bus.NewHandler(gopi.EventHandler{Name: "ikea.Event", Handler: this.EventHandlerFunc})

	// Success
	return nil
}

func (this *node) Close() error {

	// Release resources
	this.tradfri = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *node) String() string {
	str := "<" + this.Log.Name()
	str += " tradfri=" + fmt.Sprint(this.tradfri)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.Node

func (this *node) Id() string {
	return "ID TODO"
}

func (this *node) Name() string {
	return this.Unit.Log.Name()
}

func (this *node) Device(string) mutablehome.Device {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.Device

func (this *nodedevice) Id() string {
	return fmt.Sprint(this.IkeaDevice.Id())
}

func (this *nodedevice) Name() string {
	return this.IkeaDevice.Name()
}

func (this *nodedevice) Cap() []mutablehome.CapType {
	caps := make([]mutablehome.CapType, 0)

	switch this.IkeaDevice.Type() {
	case mutablehome.IKEA_DEVICE_TYPE_LIGHT:
		// If the device is a light, then power on & off are supported, and brightness
		caps = append(caps, mutablehome.CAP_POWER_OFF, mutablehome.CAP_POWER_ON, mutablehome.CAP_LIGHT_BRIGHTNESS)
	}

	return caps
}

func (this *nodedevice) Power() mutablehome.CapType {
	// If not a light, then return unknown
	if this.IkeaDevice.Type() != mutablehome.IKEA_DEVICE_TYPE_LIGHT {
		return mutablehome.CAP_NONE
	}
	// If no lights, then return unknown
	if len(this.IkeaDevice.Lights()) == 0 {
		return mutablehome.CAP_NONE
	}
	// Take first light value
	if light := this.IkeaDevice.Lights()[0]; light.Power() == true {
		return mutablehome.CAP_POWER_ON
	} else {
		return mutablehome.CAP_POWER_OFF
	}
}

func (this *nodedevice) SetPower(state mutablehome.CapType) error {
	// If not a light, then return error
	if this.IkeaDevice.Type() != mutablehome.IKEA_DEVICE_TYPE_LIGHT {
		return gopi.ErrNotImplemented.WithPrefix("Power")
	}
	// If no lights, then return error
	if len(this.IkeaDevice.Lights()) == 0 {
		return gopi.ErrNotImplemented.WithPrefix("Power")
	}
	// Set power value for all lights
	for _, light := range this.IkeaDevice.Lights() {
		switch state {
		case mutablehome.CAP_POWER_ON:
			fmt.Println("TODO: Switch on light", light)
		case mutablehome.CAP_POWER_OFF:
			fmt.Println("TODO: Switch off light", light)
		default:
			return gopi.ErrBadParameter.WithPrefix("Power")
		}
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.Evemt

func (this *node) NewEvent(t mutablehome.EventType, d mutablehome.Device) mutablehome.Event {
	return &nodeevent{t, this, d}
}

func (*nodeevent) Name() string {
	return "mutablehome.Event"
}

func (*nodeevent) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *nodeevent) Source() gopi.Unit {
	return this.Source_
}

func (this *nodeevent) Value() interface{} {
	return this.Device_
}

func (this *nodeevent) Type() mutablehome.EventType {
	return this.Type_
}

func (this *nodeevent) Node() mutablehome.Node {
	return this.Source_
}

func (this *nodeevent) Device() mutablehome.Device {
	return this.Device_
}

////////////////////////////////////////////////////////////////////////////////
// EVENT HANDLER

func (this *node) EventHandlerFunc(_ context.Context, _ gopi.App, evt gopi.Event) {
	evt_ := evt.(mutablehome.IkeaEvent)
	device_ := &nodedevice{evt_.Device()}
	switch evt_.Type() {
	case mutablehome.IKEA_EVENT_GATEWAY_CONNECTED:
		this.Emit(this.NewEvent(mutablehome.EVENT_NODE_ONLINE, nil))
	case mutablehome.IKEA_EVENT_GATEWAY_DISCONNECTED:
		this.Emit(this.NewEvent(mutablehome.EVENT_NODE_OFFLINE, nil))
	case mutablehome.IKEA_EVENT_DEVICE_ADDED:
		this.Emit(this.NewEvent(mutablehome.EVENT_DEVICE_ADDED, device_))
	case mutablehome.IKEA_EVENT_DEVICE_CHANGED:
		this.Emit(this.NewEvent(mutablehome.EVENT_DEVICE_CHANGED, device_))
	default:
		this.Log.Warn("Ignoring:", evt_)
	}
}
