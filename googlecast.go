/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type CastEventType uint

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Cast interface {
	// Return list of discovered Google Chromecast Devices
	Devices() []CastDevice

	// Connect to the control channel for a device, with timeout
	//Connect(Device, gopi.RPCFlag, time.Duration) (Channel, error)
	//Disconnect(Channel) error

	// Implements gopi.Unit
	gopi.Unit
}

type CastDevice interface {
	Id() string
	Name() string
	Model() string
	Service() string
	State() uint
}

type CastEvent interface {
	Type() CastEventType
	Device() CastDevice

	gopi.Event
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CAST_EVENT_NONE CastEventType = iota
	CAST_EVENT_ADDED
	CAST_EVENT_UPDATED
	CAST_EVENT_REMOVED
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v CastEventType) String() string {
	switch v {
	case CAST_EVENT_NONE:
		return "CAST_EVENT_NONE"
	case CAST_EVENT_ADDED:
		return "CAST_EVENT_ADDED"
	case CAST_EVENT_UPDATED:
		return "CAST_EVENT_UPDATED"
	case CAST_EVENT_REMOVED:
		return "CAST_EVENT_REMOVED"
	default:
		return "[?? Invalid CastEventType valie]"
	}
}
