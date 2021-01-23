/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package node

import (
	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Node struct {
	Gateway mutablehome.TradfriGateway
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (Node) Name() string { return "mutablehome/tradfri/node" }

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
