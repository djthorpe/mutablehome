/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Cast interface {
	// Return list of discovered Google Chromecast Devices
	Devices() []Device

	// Connect to the control channel for a device, with timeout
	//Connect(Device, gopi.RPCFlag, time.Duration) (Channel, error)
	//Disconnect(Channel) error

	// Implements gopi.Unit
	gopi.Unit
}

type Device interface {
	Id() string
	Name() string
	Model() string
	Service() string
	State() uint
}
