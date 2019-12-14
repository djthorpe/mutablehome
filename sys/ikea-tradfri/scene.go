/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type scene struct {
	Id_ uint `json:"9003"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *scene) Id() uint {
	return this.Id_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *scene) String() string {
	return fmt.Sprintf("ikea.Scene{ id=%v  }", this.Id())
}
