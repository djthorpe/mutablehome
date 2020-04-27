/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Node struct {
	Tradfri mutablehome.Ikea
}

type node struct {
	base.Unit
	tradfri mutablehome.Ikea
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
