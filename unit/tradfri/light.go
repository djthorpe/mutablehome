/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"fmt"
	"strconv"
	"time"

	"github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lightbulb struct {
	deviceId_       uint
	Power_          uint    `json:"5850"`
	ColorHex_       string  `json:"5706,omitempty"`
	ColorX_         uint16  `json:"5709,omitempty"`
	ColorY_         uint16  `json:"5710,omitempty"`
	Brightness_     uint8   `json:"5851,omitempty"`
	Temperature_    uint16  `json:"5711,omitempty"`
	TransitionTime_ float32 `json:"5712,omitempty"`
	Hue_            int     `json:"5707,omitempty"`
	Saturation_     int     `json:"5708,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this lightbulb) Power() bool {
	return this.Power_ != 0
}

func (this lightbulb) ColorHex() string {
	return this.ColorHex_
}

func (this lightbulb) ColorXY() (uint16, uint16) {
	return this.ColorX_, this.ColorY_
}

func (this lightbulb) Brightness() uint8 {
	return this.Brightness_
}

func (this lightbulb) Temperature() uint16 {
	return this.Temperature_
}

////////////////////////////////////////////////////////////////////////////////
// SET PROPERTIES

func (this lightbulb) SetPower(state bool) mutablehome.IkeaCommand {
	this.Power_ = boolToUint(state)
	return NewLightState(this.deviceId_, lightbulb{Power_: boolToUint(state)})
}

func (this lightbulb) SetBrightness(value uint8, transition time.Duration) mutablehome.IkeaCommand {
	this.Brightness_ = value
	this.TransitionTime_ = durationToTransition(transition)
	return NewLightState(this.deviceId_, lightbulb{Power_: 1, Brightness_: value, TransitionTime_: this.TransitionTime_})
}

func (this lightbulb) SetColorHex(value string, transition time.Duration) mutablehome.IkeaCommand {
	this.ColorHex_ = value
	this.TransitionTime_ = durationToTransition(transition)
	return NewLightState(this.deviceId_, lightbulb{Power_: 1, ColorHex_: value, TransitionTime_: this.TransitionTime_})
}

func (this lightbulb) SetColorXY(x uint16, y uint16, transition time.Duration) mutablehome.IkeaCommand {
	this.ColorX_, this.ColorY_ = x, y
	this.TransitionTime_ = durationToTransition(transition)
	return NewLightState(this.deviceId_, lightbulb{Power_: 1, ColorX_: x, ColorY_: y, TransitionTime_: this.TransitionTime_})
}

func (this lightbulb) SetTemperature(value uint16, transition time.Duration) mutablehome.IkeaCommand {
	this.Temperature_ = value
	this.TransitionTime_ = durationToTransition(transition)
	return NewLightState(this.deviceId_, lightbulb{Power_: 1, Temperature_: this.Temperature_, TransitionTime_: this.TransitionTime_})
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this lightbulb) String() string {
	str := "<ikea.Light" +
		" power=" + fmt.Sprint(this.Power()) +
		" brightness=" + fmt.Sprint(this.Brightness()) +
		" color_hex=" + strconv.Quote(this.ColorHex())
	x, y := this.ColorXY()
	str += " color_xy=" + fmt.Sprintf("{%v,%v}", x, y)
	return str + ">"
}
