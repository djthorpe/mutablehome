/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"io"
	"time"

	// Modules
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	TradfriDeviceType uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// TradfriGateway represents a connection to a gateway device
type TradfriGateway interface {
	gopi.Unit

	// Connect to gateway, using either IP4 or IP6
	Connect(gopi.RPCServiceRecord, gopi.RPCFlag) error

	// Disconnect from gateway
	Disconnect() error

	// Return information about the gateway
	Id() string
	Version() string

	// Return list of device, group and scene id's
	Devices() ([]uint, error)
	Groups() ([]uint, error)
	Scenes() ([]uint, error)

	/*
		// Return details of a specific device, group or scene
		Device(id uint) (TradfriDevice, error)
		Group(id uint) (TradfriGroup, error)

		// Send one or more device, group or scene commands
		Send(...TradfriCommand) error

		// Observe device
		ObserveDevice(context.Context, uint) error
	*/
}

// TradfriDevice represents a device such as a set of lights
type TradfriDevice interface {
	Id() uint
	Name() string
	Type() TradfriDeviceType
	Created() time.Time
	Updated() time.Time
	Active() bool

	Lights() []TradfriLight
}

// TradfriLight represents a single light
type TradfriLight interface {
	// Get properties
	Power() bool
	Brightness() uint8 // 00 to FE
	ColorHex() string
	ColorXY() (uint16, uint16) // 0000 to FFFF
	Temperature() uint16       // 250 to 454

	// Set properties
	SetPower(bool) TradfriCommand
	SetBrightness(uint8, time.Duration) TradfriCommand       // 01 to FE
	SetColorXY(uint16, uint16, time.Duration) TradfriCommand // 0000 to FFFF
	SetTemperature(uint16, time.Duration) TradfriCommand     // 250 to 454
	SetColorHex(string, time.Duration) TradfriCommand
}

// TradfriGroup represents a group of devices
type TradfriGroup interface {
	Id() uint
	Name() string
	Devices() []uint
}

// TradfriCommand represents a command to send to the gateway
type TradfriCommand interface {
	Path() string
	Body() (io.Reader, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	IKEA_DEFAULT_PORT = 5684
)

const (
	IKEA_DEVICE_TYPE_REMOTE       TradfriDeviceType = 0
	IKEA_DEVICE_TYPE_SLAVE_REMOTE TradfriDeviceType = 1
	IKEA_DEVICE_TYPE_LIGHT        TradfriDeviceType = 2
	IKEA_DEVICE_TYPE_PLUG         TradfriDeviceType = 3
	IKEA_DEVICE_TYPE_MOTIONSENSOR TradfriDeviceType = 4
	IKEA_DEVICE_TYPE_REPEATER     TradfriDeviceType = 6
	IKEA_DEVICE_TYPE_BLIND        TradfriDeviceType = 7
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t TradfriDeviceType) String() string {
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
		return "[?? Invalid TradfriDeviceType value]"
	}
}
