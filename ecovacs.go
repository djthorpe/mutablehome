/*
	Mutablehome Automation: Ecovacs
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"errors"

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
	Id() string

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
// GLOBALS

var (
	ErrAuthenticationError = errors.New("Authentication Error")
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

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
