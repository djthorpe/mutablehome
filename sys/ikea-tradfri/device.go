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
	"time"

	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	Name_        string `json:"9001"`
	Created_     int64  `json:"9002"`
	Updated_     int64  `json:"9020"`
	Id_          uint   `json:"9003"`
	Active_      uint   `json:"9019"`
	Type_        uint   `json:"5750"`
	NeedsUpdate_ uint   `json:"9054"`

	Metadata_ struct {
		Vendor       string `json:"0"`
		Product      string `json:"1"`
		Serial       string `json:"2"`
		Version      string `json:"3"`
		PowerSource  int    `json:"6"`
		BatteryLevel int    `json:"9"`
	} `json:"3"`

	Lights_      []lightbulb `json:"3311"` // IKEA_DEVICE_TYPE_LIGHTBULB
	MotionSensor []sensor    `json:"3300"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *device) Name() string {
	return this.Name_
}

func (this *device) Id() uint {
	return this.Id_
}

func (this *device) Type() mutablehome.IkeaDeviceType {
	return mutablehome.IkeaDeviceType(this.Type_)
}

func (this *device) Created() time.Time {
	return time.Unix(this.Created_, 0)
}

func (this *device) Updated() time.Time {
	return time.Unix(this.Updated_, 0)
}

func (this *device) Active() bool {
	return this.Active_ != 0
}

func (this *device) Vendor() string {
	return this.Metadata_.Vendor
}

func (this *device) Product() string {
	return this.Metadata_.Product
}

func (this *device) Version() string {
	return this.Metadata_.Version
}

func (this *device) Lights() []mutablehome.IkeaLight {
	lights := make([]mutablehome.IkeaLight, len(this.Lights_))
	for i, light := range this.Lights_ {
		lights[i] = light
	}
	return lights
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return fmt.Sprintf("ikea.Device{ id=%v type=%v name=%v created=%v updated=%v active=%v vendor=%v product=%v version=%v }",
		this.Id(), this.Type(), strconv.Quote(this.Name()),
		this.Created(), this.Updated(), this.Active(),
		strconv.Quote(this.Vendor()), strconv.Quote(this.Product()), strconv.Quote(this.Version()))
}
