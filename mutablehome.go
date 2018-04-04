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
type PairStatusType uint8

type Device struct {
	// Type of device (manufacturer, protocol)
	Type DeviceType `json:"type"`

	// Unique device identifier
	DeviceId uint64 `json:"id"`

	// Per-product identifier
	ProductId uint64 `json:"product"`

	// Name of the device (user-editable)
	Name string `json:"name"`

	// Location of the device (user-editable)
	Location string `json:"location"`

	// Status of the device
	PairStatus PairStatusType `json:"status"`

	// Timestamps
	TimeDiscovered time.Time `json:"discovered_timestamp,omitempty"`
	TimeUnpaired   time.Time `json:"unpaired_timestamp,omitempty"`
	TimePaired     time.Time `json:"paired_timestamp,omitempty"`
	TimeUpdated    time.Time `json:"updated_timestamp,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Devices interface {
	gopi.Driver

	// Return an existing device (or create a new device)
	Device(device_id uint64, device_type DeviceType, product_id uint64) (*Device, error)

	// Return a list of all devices
	Devices(device_type_flag DeviceType, pair_status_flag PairStatusType) []*Device

	// Pair and Unpair
	Pair(device_id uint64, device_type DeviceType) error
	Unpair(device_id uint64, device_type DeviceType) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEVICE_TYPE_ENERGENIE_CONTROL DeviceType = (1 << iota)
	DEVICE_TYPE_ENERGENIE_MONITOR
	DEVICE_TYPE_MAX             = DEVICE_TYPE_ENERGENIE_MONITOR
	DEVICE_TYPE_ANY             = DEVICE_TYPE_MAX<<1 - 1
	DEVICE_TYPE_NONE DeviceType = 0x00
)

const (
	PAIR_STATUS_DISCOVERED PairStatusType = (1 << iota)
	PAIR_STATUS_UNPAIRED
	PAIR_STATUS_PAIRED
	PAIR_STATUS_MAX                 = PAIR_STATUS_PAIRED
	PAIR_STATUS_ANY                 = PAIR_STATUS_MAX<<1 - 1
	PAIR_STATUS_NONE PairStatusType = 0x00
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
	// Now individual cases
	flags := ""
	for f := DeviceType(1); f <= DEVICE_TYPE_MAX; f <<= 1 {
		if t&f == 0 {
			continue
		}
		switch f {
		case DEVICE_TYPE_ENERGENIE_CONTROL:
			flags = flags + ",DEVICE_TYPE_ENERGENIE_CONTROL"
		case DEVICE_TYPE_ENERGENIE_MONITOR:
			flags = flags + ",DEVICE_TYPE_ENERGENIE_MONITOR"
		default:
			flags = flags + "[?? Invalid DeviceType]"
		}
	}
	return flags[1:]
}

func (s PairStatusType) String() string {
	// None and Any cases
	switch s {
	case PAIR_STATUS_NONE:
		return "PAIR_STATUS_NONE"
	case PAIR_STATUS_ANY:
		return "PAIR_STATUS_ANY"
	}
	// Now individual cases
	flags := ""
	for f := PairStatusType(1); f <= PAIR_STATUS_MAX; f <<= 1 {
		if s&f == 0 {
			continue
		}
		switch f {
		case PAIR_STATUS_DISCOVERED:
			flags = flags + ",PAIR_STATUS_DISCOVERED"
		case PAIR_STATUS_UNPAIRED:
			flags = flags + ",PAIR_STATUS_UNPAIRED"
		case PAIR_STATUS_PAIRED:
			flags = flags + ",PAIR_STATUS_PAIRED"
		default:
			flags = flags + "[?? Invalid PairStatusType]"
		}
	}
	return flags[1:]
}

func (device *Device) Hash() string {
	return Hash(device.Type, device.DeviceId)
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
