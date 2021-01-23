/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package node

import (
	"fmt"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	Type_   mutablehome.EventType
	Source_ mutablehome.Node
	Device_ mutablehome.Device
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *node) NewGatewayEvent(t mutablehome.EventType) mutablehome.Event {
	return &event{t, this, nil}
}

func (this *node) NewDeviceEvent(t mutablehome.EventType, d mutablehome.Device) mutablehome.Event {
	return &event{t, this, d}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (*event) Name() string {
	return "mutablehome.Event"
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

func (this *event) Type() mutablehome.EventType {
	return this.Type_
}

func (this *event) Node() mutablehome.Node {
	return this.Source_
}

func (this *event) Device() mutablehome.Device {
	return this.Device_
}

func (this *event) Traits() []mutablehome.TraitType {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<" + this.Name()
	str += " type=" + fmt.Sprint(this.Type_)
	if this.Device_ != nil {
		str += " device=" + fmt.Sprint(this.Device_)
	}
	return str + ">"
}
