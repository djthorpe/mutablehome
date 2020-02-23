/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	Source_ gopi.Unit
	Type_   mutablehome.IkeaEventType
	Device_ mutablehome.IkeaDevice
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewEvent(source gopi.Unit, type_ mutablehome.IkeaEventType, device mutablehome.IkeaDevice) mutablehome.IkeaEvent {
	return &event{source, type_, device}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*event) Name() string {
	return "ikea.Event"
}

func (*event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *event) Source() gopi.Unit {
	return this.Source_
}

func (this *event) Value() interface{} {
	return this.Device_
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION ikea.Event

func (this *event) Type() mutablehome.IkeaEventType {
	return this.Type_
}

func (this *event) Device() mutablehome.IkeaDevice {
	return this.Device_
}
