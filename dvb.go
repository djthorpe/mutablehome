/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks
	gopi2 "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type DVBTable interface {

	// Properties returns an array of DVB Properties
	Properties() []DVBProperties

	gopi2.Unit
}

type DVBProperties interface {
}
