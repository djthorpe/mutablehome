/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks

	"context"

	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type CastEventType uint

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Cast interface {
	// Return list of discovered Google Chromecast Devices
	Devices(context.Context) ([]CastDevice, error)

	// Connect to the control channel for a device
	Connect(CastDevice, gopi.RPCFlag) error

	// Disconnect from the device
	Disconnect(CastDevice) error

	// Implements gopi.Unit
	gopi.Unit
}

type CastDevice interface {
	Id() string
	Name() string
	Model() string
	Service() string
	State() uint

	// Volume
	Volume() CastVolume
	SetVolume(level float32) error
	SetMute(mute bool) error

	// Application
	App() CastApp
	LaunchAppWithId(string) error

	// Play, pause and stop
	SetPlay(bool) error  // Play or stop
	SetPause(bool) error // Pause or play

	// Load Media by URL
	LoadURL(url string, autoplay bool) error
}

type CastEvent interface {
	Type() CastEventType
	Device() CastDevice

	gopi.Event
}

type CastVolume interface {
	Level() float32
	Muted() bool
}

type CastApp interface {
	ID() string
	Name() string
	Status() string
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
