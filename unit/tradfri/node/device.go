/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package node

import (
	"fmt"

	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	id   uint
	seen bool
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewDevice(id uint) *device {
	this := new(device)
	this.id = id
	return this
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *device) Id() string {
	return fmt.Sprint(this.id)
}

func (this *device) Name() string {
	return "TODO"
}

func (this *device) Traits() []mutablehome.TraitType {
	caps := make([]mutablehome.TraitType, 0)
	return caps
}
