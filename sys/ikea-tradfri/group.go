/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type group struct {
	Id_      uint   `json:"9003"`
	Name_    string `json:"9001"`
	Created_ int64  `json:"9002"`
	Content  struct {
		Devices struct {
			Devices []uint `json:"9003"`
		} `json:"15002"`
	} `json:"9018"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *group) Id() uint {
	return this.Id_
}

func (this *group) Name() string {
	return this.Name_
}

func (this *group) Devices() []uint {
	return this.Content.Devices.Devices
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *group) String() string {
	return fmt.Sprintf("ikea.Group{ id=%v name=%v devices=%v }", this.Id(), strconv.Quote(this.Name()), this.Devices())
}
