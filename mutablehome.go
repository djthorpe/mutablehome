/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package mutablehome

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DeviceType uint64
type PairStatus uint8

type Device struct {
	// Type of device (manufacturer, protocol)
	Type DeviceType

	// Unique device identifier
	DeviceId uint64

	// Per-product identifier
	ProductId uint64

	// Name of the device (user-editable)
	Name string

	// Location of the device (user-editable)
	Location string

	// Status of the device
	Paired PairStatus

	// Timestamps
	TimeDiscovered time.Time
	TimeUnpaired   time.Time
	TimePaired     time.Time
	TimeUpdated    time.Time
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Devices interface {
	gopi.Driver

	// Return Hashkey for a device or empty string
	//HashForDevice(*Device) string

	// Return an existing device (or create a new device)
	Device(device_id uint64, device_type DeviceType, product_id uint64) (*Device, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEVICE_TYPE_NONE              DeviceType = 0x00
	DEVICE_TYPE_ENERGENIE_CONTROL DeviceType = (1 << iota)
	DEVICE_TYPE_ENERGENIE_MONITOR DeviceType = (1 << iota)
	DEVICE_TYPE_ANY               DeviceType = (1 << iota) - 1
	DEVICE_TYPE_MAX               DeviceType = DEVICE_TYPE_ENERGENIE_MONITOR
)

const (
	PAIR_STATUS_NONE       PairStatus = 0x00
	PAIR_STATUS_DISCOVERED PairStatus = (1 << iota)
	PAIR_STATUS_UNPAIRED   PairStatus = (1 << iota)
	PAIR_STATUS_PAIRED     PairStatus = (1 << iota)
	PAIR_STATUS_ANY        PairStatus = (1 << iota) - 1
	PAIR_STATUS_MAX        PairStatus = PAIR_STATUS_PAIRED
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t DeviceType) String() string {
	// None and Any cases
	switch t {
	case DEVICE_TYPE_NONE:
		return "DEVICE_TYPE_NONE"
	case DEVICE_TYPE_ANY:
		return "DEVICE_TYPE_ANY"
	}
	// TODO
	return "DEVICE_TYPE_OTHER"
}

func (s PairStatus) String() string {
	// None and Any cases
	switch s {
	case PAIR_STATUS_NONE:
		return "PAIR_STATUS_NONE"
	case PAIR_STATUS_ANY:
		return "PAIR_STATUS_ANY"
	}
	// TODO
	return "PAIR_STATUS_OTHER"
}

////////////////////////////////////////////////////////////////////////////////
// HASH KEY

// Hash returns unique key as a string for a particular device or empty string
// on error
func Hash(device_type DeviceType, device_id uint64) string {
	switch device_type {
	case DEVICE_TYPE_ENERGENIE_CONTROL:
		return hashEnergenieControl(device_id)
	case DEVICE_TYPE_ENERGENIE_MONITOR:
		return hashEnergenieMonitor(device_id)
	default:
		return ""
	}
}

func hashEnergenieControl(device_id uint64) string {
	// 20-bit identifier
	if device_id&0xFFFFF != device_id {
		return ""
	} else {
		return fmt.Sprintf("EC:%05X", device_id)
	}
}

func hashEnergenieMonitor(device_id uint64) string {
	// 24-bit identifier
	if device_id&0xFFFFFF != device_id {
		return ""
	} else {
		return fmt.Sprintf("EM:%06X", device_id)
	}
}
