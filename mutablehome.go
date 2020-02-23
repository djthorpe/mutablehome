/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"time"

	// Frameworks
	gopi2 "github.com/djthorpe/gopi/v2"
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

////////////////////////////////////////////////////////////////////////////////
// IKEA TRADFRI

type IkeaDeviceType uint

type Tradfri interface {
	/*
		// Connect to gateway, using either IP4 or IP6
		Connect(gopi.RPCServiceRecord, gopi.RPCFlag) error

		// Return list of devices, groups and scenes
		Devices() ([]uint, error)
		Groups() ([]uint, error)
		Scenes() ([]uint, error)

		// Return details of a specific device, group or scene
		Device(id uint) (IkeaDevice, error)
		Group(id uint) (IkeaGroup, error)
		Scene(id uint) (IkeaScene, error)

		// Observe devices
		ObserveDevice(context.Context, uint) error
	*/

	// Implements Unit
	gopi2.Unit
}

type IkeaDevice interface {
	Id() uint
	Name() string
	Type() IkeaDeviceType
	Created() time.Time
	Updated() time.Time
	Active() bool

	Lights() []IkeaLight
}

type IkeaLight interface {
	Power() bool
	ColorHex() string
}

type IkeaGroup interface {
	Id() uint
	Name() string
	Devices() []uint
}

type IkeaScene interface {
	Id() uint
}

const (
	IKEA_DEVICE_TYPE_REMOTE       IkeaDeviceType = 0
	IKEA_DEVICE_TYPE_SLAVE_REMOTE IkeaDeviceType = 1
	IKEA_DEVICE_TYPE_LIGHT        IkeaDeviceType = 2
	IKEA_DEVICE_TYPE_PLUG         IkeaDeviceType = 3
	IKEA_DEVICE_TYPE_MOTIONSENSOR IkeaDeviceType = 4
	IKEA_DEVICE_TYPE_REPEATER     IkeaDeviceType = 6
	IKEA_DEVICE_TYPE_BLIND        IkeaDeviceType = 7
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t IkeaDeviceType) String() string {
	switch t {
	case IKEA_DEVICE_TYPE_REMOTE:
		return "IKEA_DEVICE_TYPE_REMOTE"
	case IKEA_DEVICE_TYPE_SLAVE_REMOTE:
		return "IKEA_DEVICE_TYPE_SLAVE_REMOTE"
	case IKEA_DEVICE_TYPE_LIGHT:
		return "IKEA_DEVICE_TYPE_LIGHT"
	case IKEA_DEVICE_TYPE_PLUG:
		return "IKEA_DEVICE_TYPE_PLUG"
	case IKEA_DEVICE_TYPE_MOTIONSENSOR:
		return "IKEA_DEVICE_TYPE_MOTIONSENSOR"
	case IKEA_DEVICE_TYPE_REPEATER:
		return "IKEA_DEVICE_TYPE_REPEATER"
	case IKEA_DEVICE_TYPE_BLIND:
		return "IKEA_DEVICE_TYPE_BLIND"
	default:
		return "[?? Invalid IkeaDeviceType value]"
	}
}
