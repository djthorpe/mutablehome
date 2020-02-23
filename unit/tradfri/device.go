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

	// Frameworks
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

	Lights_ []lightbulb `json:"3311"` // IKEA_DEVICE_TYPE_LIGHTBULB
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func NewDevice() *device {
	return new(device)
}

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
		light.deviceId_ = this.Id()
		lights[i] = light
	}
	return lights
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	str := "<ikea.Device" +
		" id=" + fmt.Sprint(this.Id()) +
		" type=" + fmt.Sprint(this.Type()) +
		" name=" + strconv.Quote(this.Name()) +
		" active=" + fmt.Sprint(this.Active())

	switch this.Type() {
	case mutablehome.IKEA_DEVICE_TYPE_LIGHT:
		str += " lights=" + fmt.Sprint(this.Lights())
	}

	if created := this.Created(); created.IsZero() == false {
		str += " created=" + created.Format(time.RFC3339)
	}
	if updated := this.Updated(); updated.IsZero() == false {
		str += " updated=" + updated.Format(time.RFC3339)
	}
	if vendor := this.Vendor(); vendor != "" {
		str += " vendor=" + strconv.Quote(vendor)
	}
	if product := this.Product(); product != "" {
		str += " product=" + strconv.Quote(product)
	}
	if version := this.Version(); version != "" {
		str += " version=" + strconv.Quote(version)
	}
	return str + ">"
}
