/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"context"
	"io"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// IKEA TRADFRI

type IkeaDeviceType uint

type Ikea interface {
	// Connect to gateway, using either IP4 or IP6
	Connect(gopi.RPCServiceRecord, gopi.RPCFlag) error

	// Disconnect from gateway
	Disconnect() error

	// Return list of device id's
	Devices() ([]uint, error)
	Groups() ([]uint, error)
	Scenes() ([]uint, error)

	// Return details of a specific device, group or scene
	Device(id uint) (IkeaDevice, error)
	Group(id uint) (IkeaGroup, error)

	// Send one or more device, group or scene commands
	Send(...IkeaCommand) error

	// Observe devices
	ObserveDevice(context.Context, uint) error

	// Implements Unit
	gopi.Unit
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
	// Get properties
	Power() bool
	Brightness() uint8 // 00 to FE
	ColorHex() string
	ColorXY() (uint16, uint16) // 0000 to FFFF
	Temperature() uint16       // 250 to 454

	// Set properties
	SetPower(bool) IkeaCommand
	SetBrightness(uint8, time.Duration) IkeaCommand       // 01 to FE
	SetColorXY(uint16, uint16, time.Duration) IkeaCommand // 0000 to FFFF
	SetTemperature(uint16, time.Duration) IkeaCommand     // 250 to 454
	SetColorHex(string, time.Duration) IkeaCommand
}

type IkeaGroup interface {
	Id() uint
	Name() string
	Devices() []uint
}

type IkeaCommand interface {
	Path() string
	Body() (io.Reader, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	IKEA_DEFAULT_PORT = 5684
)

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
