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

type lightbulb struct {
	Power_          uint    `json:"5850,omitempty"`
	ColorHex_       string  `json:"5706,omitempty"`
	ColorX_         int16   `json:"5709,omitempty"`
	ColorY_         int16   `json:"5710,omitempty"`
	Hue_            int     `json:"5707,omitempty"`
	Saturation_     int     `json:"5708,omitempty"`
	Temperature_    int     `json:"5711,omitempty"`
	Brightness_     int     `json:"5851,omitempty"`
	TransitionTime_ float32 `json:"5712,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this lightbulb) Power() bool {
	return this.Power_ != 0
}

func (this lightbulb) ColorHex() string {
	return this.ColorHex_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this lightbulb) String() string {
	return fmt.Sprintf("ikea.Light{ power=%v color=%v }", this.Power(), strconv.Quote(this.ColorHex()))
}
