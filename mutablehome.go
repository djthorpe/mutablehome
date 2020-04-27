/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PowerState uint

////////////////////////////////////////////////////////////////////////////////
// NODE

// Node is a collection of devices
type Node interface {
	NodeName() string
	Devices() []Device
}

// Device is a device which can be controlled
type Device interface {
	Id() string   // Globally unique ID for the device
	Name() string // Name of the device
}

// CapPower represents a device which can be switched on, off or toggled
type CapPower interface {
	Device

	State() PowerState
	On() error
	Off() error
}

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
