/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"context"
	"time"

	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EventType uint
	TraitType uint
)

////////////////////////////////////////////////////////////////////////////////
// NODE

// Node is a collection of devices which can be observed or controlled
// and accessed by a unique ID
type Node interface {
	gopi.PubSub

	Id() string           // Unique Id for the node
	Name() string         // Textual description of the node
	Device(string) Device // Return device with Id
}

// Device is a device which can be observed or controlled
type Device interface {
	Id() string          // Unique ID for the device
	Name() string        // Name of the device
	Traits() []TraitType // Capabilities for the device
}

// PowerTrait represents a device which can be switched on, off or toggled
type PowerTrait interface {
	Device

	Power() TraitType         // Return ON, OFF or STANDBY or NONE if unknown
	SetPower(TraitType) error // Set ON, OFF, STANDBY or TOGGLE
}

// LightTrait represents a device which can have brightness or hue set
type LightTrait interface {
	Device

	Brightness() float32                        // Return brightness between 0.0 and 1.0
	SetBrightness(float32, time.Duration) error // Set brightness between 0.0 and 1.0 and a transition time
}

// Event is emitted when a device changes or node is online or offline
// of type mutablehome.Event
type Event interface {
	gopi.Event

	Type() EventType
	Node() Node
	Device() Device
	Traits() []TraitType
}

// NodeStub represents a connection to a remote mutablehome node
type NodeStub interface {
	gopi.RPCClientStub

	// Ping returns without error if the remote service is running
	Ping(context.Context) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TRAIT_NONE TraitType = iota
	TRAIT_POWER_ON
	TRAIT_POWER_OFF
	TRAIT_POWER_STANDBY
	TRAIT_POWER_TOGGLE
	TRAIT_LIGHT_BRIGHTNESS
	TRAIT_LIGHT_TEMPERATURE
	TRAIT_LIGHT_COLOR
	TRAIT_LIGHT_TRANSITION
)

const (
	EVENT_NONE EventType = iota
	EVENT_NODE_ONLINE
	EVENT_NODE_OFFLINE
	EVENT_DEVICE_ADDED
	EVENT_DEVICE_REMOVED
	EVENT_DEVICE_METADATA_CHANGED
	EVENT_DEVICE_TRAIT_CHANGED
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v EventType) String() string {
	switch v {
	case EVENT_NONE:
		return "EVENT_NONE"
	case EVENT_NODE_ONLINE:
		return "EVENT_NODE_ONLINE"
	case EVENT_NODE_OFFLINE:
		return "EVENT_NODE_OFFLINE"
	case EVENT_DEVICE_ADDED:
		return "EVENT_DEVICE_ADDED"
	case EVENT_DEVICE_REMOVED:
		return "EVENT_DEVICE_REMOVED"
	case EVENT_DEVICE_METADATA_CHANGED:
		return "EVENT_DEVICE_METADATA_CHANGED"
	case EVENT_DEVICE_TRAIT_CHANGED:
		return "EVENT_DEVICE_TRAIT_CHANGED"
	default:
		return "[?? Invalid EventType value]"
	}
}

func (v TraitType) String() string {
	switch v {
	case TRAIT_NONE:
		return "TRAIT_NONE"
	case TRAIT_POWER_ON:
		return "TRAIT_POWER_ON"
	case TRAIT_POWER_OFF:
		return "TRAIT_POWER_OFF"
	case TRAIT_POWER_STANDBY:
		return "TRAIT_POWER_STANDBY"
	case TRAIT_POWER_TOGGLE:
		return "TRAIT_POWER_TOGGLE"
	case TRAIT_LIGHT_BRIGHTNESS:
		return "TRAIT_LIGHT_BRIGHTNESS"
	case TRAIT_LIGHT_TEMPERATURE:
		return "TRAIT_LIGHT_TEMPERATURE"
	case TRAIT_LIGHT_COLOR:
		return "TRAIT_LIGHT_COLOR"
	case TRAIT_LIGHT_TRANSITION:
		return "TRAIT_LIGHT_TRANSITION"
	default:
		return "[?? Invalid TraitType value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// INFLUXDB (old code)

type InfluxDB interface {
	// Create a new resultset
	NewResultSet(tags map[string]string) InfluxRS

	// Write rows to the database
	Write(InfluxRS) error
}

type InfluxRS interface {
	// Remove all rows
	RemoveAll()

	// Add a new row of metrics for a measurement name
	Add(string, map[string]interface{}) error

	// Add a new row of metrics using timestamp for a measurement name
	AddTS(string, map[string]interface{}, time.Time) error
}
