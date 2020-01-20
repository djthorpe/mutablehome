package ecovacs

import (
	"fmt"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	type_   home.EcovacsEventType
	source  home.Ecovacs
	device  home.EvovacsDevice
	reqId   string
	data    []byte
	message *message
}

////////////////////////////////////////////////////////////////////////////////
// NEW EVENT

func NewEvent(source home.Ecovacs, device home.EvovacsDevice, reqId string, msg *message, data []byte) *event {
	evt := &event{
		type_:   home.ECOVACS_EVENT_NONE,
		source:  source,
		device:  device,
		reqId:   reqId,
		message: msg,
		data:    data,
	}
	if msg.Control.Battery.Power > 0 {
		evt.type_ = home.ECOVACS_EVENT_BATTERYLEVEL
	} else if msg.Control.Clean.Type != "" {
		evt.type_ = home.ECOVACS_EVENT_CLEANSTATE
	} else if msg.Control.Charge.Type != "" {
		evt.type_ = home.ECOVACS_EVENT_CHARGESTATE
	} else if msg.Control.Type != "" {
		evt.type_ = home.ECOVACS_EVENT_LIFESPAN
	} else if msg.Control.ErrorNo != 0 {
		evt.type_ = home.ECOVACS_EVENT_ERROR
	} else if msg.Control.Version.Name != "" {
		evt.type_ = home.ECOVACS_EVENT_VERSION
	}
	return evt
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*event) Name() string {
	return "mutablehome.EcovacsEvent"
}

func (*event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *event) Source() gopi.Unit {
	return this.source
}

func (this *event) Value() interface{} {
	switch this.type_ {
	case home.ECOVACS_EVENT_BATTERYLEVEL:
		return this.BatteryLevel()
	case home.ECOVACS_EVENT_CLEANSTATE:
		mode, suction := this.CleanState()
		return []interface{}{mode, suction}
	case home.ECOVACS_EVENT_CHARGESTATE:
		return this.ChargeState()
	case home.ECOVACS_EVENT_LIFESPAN:
		part, val, total := this.LifeSpan()
		return []interface{}{part, val, total}
	case home.ECOVACS_EVENT_VERSION:
		return this.Version()
	case home.ECOVACS_EVENT_ERROR:
		errNum, msg := this.Error()
		return []interface{}{errNum, msg}
	default:
		return this.data
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.EcovacsEvent

func (this *event) Type() home.EcovacsEventType {
	return this.type_
}

func (this *event) Device() home.EvovacsDevice {
	return this.device
}

func (this *event) RequestId() string {
	return this.reqId
}

func (this *event) BatteryLevel() uint {
	return this.message.Control.Battery.Power
}

func (this *event) Error() (uint, string) {
	return this.message.Control.ErrorNo, this.ErrorMsg()
}

func (this *event) ErrorMsg() string {
	if this.message.Control.Error != "" {
		return this.message.Control.Error
	} else {
		switch this.message.Control.ErrorNo {
		case 100:
			return "NoError"
		case 101:
			return "BatteryLow"
		case 102:
			return "HostHang"
		case 103:
			return "WheelAbnormal"
		case 104:
			return "DownSensorAbnormal"
		case 110:
			return "NoDustBox"
		default:
			return fmt.Sprintf("Error%03d", this.message.Control.ErrorNo)
		}
	}
}

func (this *event) ChargeState() string {
	return strings.ToLower(this.message.Control.Charge.Type)
}

func (this *event) CleanState() (home.EcovacsCleanMode, home.EcovacsCleanSuction) {
	return home.EcovacsCleanMode(strings.ToLower(this.message.Control.Clean.Type)), home.EcovacsCleanSuction(strings.ToLower(this.message.Control.Clean.Speed))
}

func (this *event) LifeSpan() (home.EcovacsPart, uint, uint) {
	return home.EcovacsPart(this.message.Control.Type), this.message.Control.Val, this.message.Control.Total
}

func (this *event) Version() string {
	return this.message.Control.Version.Value
}

////////////////////////////////////////////////////////////////////////////////
// EQUALS

func (this *event) Equals(other home.EcovacsEvent) bool {
	// Never matches if types are different
	if this.Type() != other.Type() {
		return false
	}
	// Try and match values which are either uint or string
	value := this.Value()
	switch value.(type) {
	case string:
		if other.Value().(string) == value {
			return true
		} else {
			return false
		}
	case uint:
		if other.Value().(uint) == value {
			return true
		} else {
			return false
		}
	}
	// For clean state, we match the mode and suction
	switch this.Type() {
	case home.ECOVACS_EVENT_CLEANSTATE:
		thisValues := this.Value().([]interface{})
		otherValues := other.Value().([]interface{})
		if len(thisValues) != 2 || len(otherValues) != 2 {
			return false
		}
		if thisValues[0] != otherValues[0] {
			return false
		}
		if thisValues[1] != otherValues[1] {
			return false
		}
		return true
	case home.ECOVACS_EVENT_LIFESPAN:
		thisValues := this.Value().([]interface{})
		otherValues := other.Value().([]interface{})
		if len(thisValues) != 3 || len(otherValues) != 3 {
			return false
		}
		if thisValues[0] != otherValues[0] {
			return false
		}
		if thisValues[1] != otherValues[1] {
			return false
		}
		if thisValues[2] != otherValues[2] {
			return false
		}
		return true
	}

	// In default case, indicate no equals
	return false
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<" + this.Name() +
		" type=" + fmt.Sprint(this.Type()) +
		" request_id=" + this.RequestId()
	if this.Device() != nil {
		str += " device=" + fmt.Sprint(this.Device())
	}
	if this.type_ == home.ECOVACS_EVENT_NONE {
		str += " data=" + strconv.Quote(string(this.data))
	} else {
		str += " value=" + fmt.Sprint(this.Value())
	}
	str += ">"
	return str
}
