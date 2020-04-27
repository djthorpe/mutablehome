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

func (this *node) Devices() []mutablehome.Device {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// EVENT HANDLER

func (this *node) EventHandlerFunc(_ context.Context, _ gopi.App, evt gopi.Event) {
	evt_ := evt.(mutablehome.IkeaEvent)
	switch evt_.Type() {
	case mutablehome.IKEA_EVENT_GATEWAY_CONNECTED:
		this.Log.Info("Gateway connected")
	case mutablehome.IKEA_EVENT_GATEWAY_DISCONNECTED:
		this.Log.Info("Gateway disconnected")
	case mutablehome.IKEA_EVENT_DEVICE_ADDED:
		this.Log.Info("Device added", evt_.Device())
	case mutablehome.IKEA_EVENT_DEVICE_CHANGED:
		this.Log.Info("Device changed", evt_.Device())
	default:
		this.Log.Warn("Ignoring:", evt_)
	}
}
