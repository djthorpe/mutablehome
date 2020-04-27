/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"time"

	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EventType uint
	CapType   uint
)

////////////////////////////////////////////////////////////////////////////////
// NODE

// Node is a collection of devices
type Node interface {
	gopi.PubSub

	Id() string           // Unique Id for the node
	Name() string         // Textual description of the node
	Device(string) Device // Return device with Id
}

// Device is a device which can be controlled
type Device interface {
	Id() string     // Globally unique ID for the device
	Name() string   // Name of the device
	Cap() []CapType // Capabilities for the device
}

// PowerCapability represents a device which can be switched on, off or toggled
type PowerCapability interface {
	Device

	PowerState() CapType // Return ON, OFF or STANDBY
	PowerOn() error
	PowerOff() error
	PowerStandby() error
	PowerToggle() error
}

// Event is emitted when a device changes or node is online or offline
// of type mutablehome.Event
type Event interface {
	gopi.Event

	Type() EventType
	Node() Node
	Device() Device
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CAP_NONE CapType = iota
	CAP_POWER_ON
	CAP_POWER_OFF
	CAP_POWER_STANDBY
	CAP_POWER_TOGGLE
)

const (
	EVENT_NONE EventType = iota
	EVENT_NODE_ONLINE
	EVENT_NODE_OFFLINE
	EVENT_DEVICE_ADDED
	EVENT_DEVICE_REMOVED
	EVENT_DEVICE_UPDATED
)

////////////////////////////////////////////////////////////////////////////////
// INFLUXDB

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
