/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	source_ gopi.Unit
	type_   mutablehome.CastEventType
	device_ mutablehome.CastDevice
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewAddedEvent(source gopi.Unit, device mutablehome.CastDevice) gopi.Event {
	return &event{source, mutablehome.CAST_EVENT_ADDED, device}
}

func NewUpdatedEvent(source gopi.Unit, device mutablehome.CastDevice) gopi.Event {
	return &event{source, mutablehome.CAST_EVENT_UPDATED, device}
}

func NewRemovedEvent(source gopi.Unit, device mutablehome.CastDevice) gopi.Event {
	return &event{source, mutablehome.CAST_EVENT_REMOVED, device}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (*event) Name() string {
	return "cast.Event"
}

func (*event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *event) Source() gopi.Unit {
	return this.source_
}

func (this *event) Type() mutablehome.CastEventType {
	return this.type_
}

func (this *event) Value() interface{} {
	return this.device_
}

func (this *event) Device() mutablehome.CastDevice {
	return this.device_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	return fmt.Sprintf("<%s type=%v device=%v>", this.Name(), this.type_, this.device_)
}
