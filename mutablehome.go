/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"errors"
	"time"

	// Frameworks

	gopi2 "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// ECOVACS DEEBOT

type (
	EcovacsEventType    uint
	EcovacsPart         string
	EcovacsCleanMode    string
	EcovacsCleanSuction string
)

type Ecovacs interface {
	gopi2.Unit

	// Authenticate
	Authenticate() error

	// Devices
	Devices() ([]EvovacsDevice, error)

	// Connect to a device to start reading messages
	Connect(EvovacsDevice) error

	// Disconnect from a device
	Disconnect(EvovacsDevice) error
}

type EvovacsDevice interface {
	// Return device properties
	Address() string
	Nickname() string

	// Fetch information from device, returns ReqId for the request
	GetBatteryInfo() (string, error)
	GetLifeSpan(EcovacsPart) (string, error)
	GetChargeState() (string, error)
	GetCleanState() (string, error)
	GetVersion() (string, error)

	// Command the device
	Clean(EcovacsCleanMode, EcovacsCleanSuction) (string, error)
	Charge() (string, error)
}

type EcovacsEvent interface {
	Type() EcovacsEventType
	Device() EvovacsDevice
	Id() string

	gopi2.Event
}

const (
	ECOVACS_EVENT_NONE EcovacsEventType = iota
	ECOVACS_EVENT_BATTERYLEVEL
	ECOVACS_EVENT_CLEANSTATE
	ECOVACS_EVENT_CHARGESTATE
	ECOVACS_EVENT_LIFESPAN
	ECOVACS_EVENT_TIME
	ECOVACS_EVENT_VERSION
	ECOVACS_EVENT_LOG
	ECOVACS_EVENT_ERROR
)

const (
	ECOVACS_PART_BRUSH      EcovacsPart = "Brush"
	ECOVACS_PART_SIDEBRUSH  EcovacsPart = "SideBrush"
	ECOVACS_PART_DUSTFILTER EcovacsPart = "DustCaseHeap"
)

const (
	ECOVACS_CLEAN_STOP   EcovacsCleanMode = "stop"
	ECOVACS_CLEAN_AUTO   EcovacsCleanMode = "auto"
	ECOVACS_CLEAN_BORDER EcovacsCleanMode = "border"
	ECOVACS_CLEAN_SPOT   EcovacsCleanMode = "spot"
	ECOVACS_CLEAN_ROOM   EcovacsCleanMode = "singleroom"
)

const (
	ECOVACS_SUCTION_STANDARD EcovacsCleanSuction = "standard"
	ECOVACS_SUCTION_STRONG   EcovacsCleanSuction = "strong"
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
// GLOBALS

var (
	ErrAuthenticationError = errors.New("Authentication Error")
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

func (v EcovacsEventType) String() string {
	switch v {
	case ECOVACS_EVENT_NONE:
		return "ECOVACS_EVENT_NONE"
	case ECOVACS_EVENT_BATTERYLEVEL:
		return "ECOVACS_EVENT_BATTERYLEVEL"
	case ECOVACS_EVENT_CLEANSTATE:
		return "ECOVACS_EVENT_CLEANSTATE"
	case ECOVACS_EVENT_CHARGESTATE:
		return "ECOVACS_EVENT_CHARGESTATE"
	case ECOVACS_EVENT_LIFESPAN:
		return "ECOVACS_EVENT_LIFESPAN"
	case ECOVACS_EVENT_TIME:
		return "ECOVACS_EVENT_TIME"
	case ECOVACS_EVENT_VERSION:
		return "ECOVACS_EVENT_VERSION"
	case ECOVACS_EVENT_LOG:
		return "ECOVACS_EVENT_LOG"
	case ECOVACS_EVENT_ERROR:
		return "ECOVACS_EVENT_ERROR"
	default:
		return "[?? Invalid EcovacsEventType value]"
	}
}
