package ecovacs

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EcovacsEvent struct {
	source  home.Ecovacs
	device  home.EvovacsDevice
	message *XMPPMessage
}

////////////////////////////////////////////////////////////////////////////////
// NEW EVENT

func NewEvent(source home.Ecovacs, device home.EvovacsDevice, message *XMPPMessage) *EcovacsEvent {
	return &EcovacsEvent{
		source:  source,
		device:  device,
		message: message,
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*EcovacsEvent) Name() string {
	return "mutablehome.EcovacsEvent"
}

func (*EcovacsEvent) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *EcovacsEvent) Source() gopi.Unit {
	return this.source
}

func (this *EcovacsEvent) Value() interface{} {
	return this.message.Value()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.EcovacsEvent

func (this *EcovacsEvent) Type() home.EcovacsEventType {
	return this.message.Type()
}

func (this *EcovacsEvent) Device() home.EvovacsDevice {
	return this.device
}

func (this *EcovacsEvent) RequestId() string {
	return this.message.Id()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *EcovacsEvent) String() string {
	str := "<" + this.Name()
	if this.RequestId() != "" {
		str += " request_id=" + this.RequestId()
	}
	if this.Device() != nil {
		str += " device=" + fmt.Sprint(this.Device())
	}
	if this.message != nil {
		str += " message=" + fmt.Sprint(this.message)
	}
	str += ">"
	return str
}
