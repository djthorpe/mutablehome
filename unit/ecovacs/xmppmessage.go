package ecovacs

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	// Frameworks
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// XMPPMessage encapsulates all the message information which can be returned from
// DeeBot
type XMPPMessage struct {
	XMLName xml.Name `xml:"query"`

	// Data structure
	Control struct {
		Id      string `xml:"id,attr"`
		Ret     string `xml:"ret,attr"`
		ErrorNo uint   `xml:"errno,attr"`
		Error   string `xml:"error,attr"`
		Type    string `xml:"type,attr"`
		Val     uint   `xml:"val,attr"`
		Total   uint   `xml:"total,attr"`
		Battery struct {
			Power uint `xml:"power,attr"`
		} `xml:"battery"`
		Charge struct {
			Type string `xml:"type,attr"`
		} `xml:"charge"`
		Clean struct {
			Type  string `xml:"type,attr"`
			Speed string `xml:"speed,attr"`
		} `xml:"clean"`
		Version struct {
			Name  string `xml:"name,attr"`
			Value string `xml:",chardata"`
		} `xml:"ver"`
	} `xml:"ctl"`

	// The raw data
	data []byte

	// The message ID
	id string

	// Cache type
	messageType home.EcovacsEventType
}

////////////////////////////////////////////////////////////////////////////////
// NEW

// NewXMPPMessage creates a new message
func NewXMPPMessage(data []byte, id string) (*XMPPMessage, error) {
	message := new(XMPPMessage)
	message.data = data
	message.id = id
	if err := xml.Unmarshal(data, &message); err != nil {
		return message, err
	} else {
		return message, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

// Type returns the type of message which the message represents
func (this *XMPPMessage) Type() home.EcovacsEventType {
	// Returned cached version
	if this.messageType != home.ECOVACS_EVENT_NONE {
		return this.messageType
	}
	// Set cached version
	switch {
	case this.Control.Battery.Power > 0:
		this.messageType = home.ECOVACS_EVENT_BATTERYLEVEL
	case this.Control.Clean.Type != "":
		this.messageType = home.ECOVACS_EVENT_CLEANSTATE
	case this.Control.Charge.Type != "":
		this.messageType = home.ECOVACS_EVENT_CHARGESTATE
	case this.Control.Type != "":
		this.messageType = home.ECOVACS_EVENT_LIFESPAN
	case this.Control.ErrorNo != 0:
		this.messageType = home.ECOVACS_EVENT_ERROR
	case this.Control.Version.Name != "":
		this.messageType = home.ECOVACS_EVENT_VERSION
	}
	// Return cached
	return this.messageType
}

// Value returns the information contained within the message
func (this *XMPPMessage) Value() interface{} {
	switch this.Type() {
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

// ReqId returns the ID for the message
func (this *XMPPMessage) Id() string {
	return this.id
}

func (this *XMPPMessage) BatteryLevel() uint {
	return this.Control.Battery.Power
}

func (this *XMPPMessage) Error() (uint, string) {
	return this.Control.ErrorNo, this.ErrorMsg()
}

func (this *XMPPMessage) ErrorMsg() string {
	if this.Control.Error != "" {
		return this.Control.Error
	} else {
		switch this.Control.ErrorNo {
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
			return fmt.Sprintf("Error%03d", this.Control.ErrorNo)
		}
	}
}

func (this *XMPPMessage) ChargeState() string {
	return strings.ToLower(this.Control.Charge.Type)
}

func (this *XMPPMessage) CleanState() (home.EcovacsCleanMode, home.EcovacsCleanSuction) {
	return home.EcovacsCleanMode(strings.ToLower(this.Control.Clean.Type)), home.EcovacsCleanSuction(strings.ToLower(this.Control.Clean.Speed))
}

func (this *XMPPMessage) LifeSpan() (home.EcovacsPart, uint, uint) {
	return home.EcovacsPart(this.Control.Type), this.Control.Val, this.Control.Total
}

func (this *XMPPMessage) Version() string {
	return this.Control.Version.Value
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *XMPPMessage) String() string {
	switch this.Type() {
	case home.ECOVACS_EVENT_NONE:
		return "<XMPPMessage" +
			" id=" + this.id +
			" data=" + strconv.Quote(string(this.data)) +
			">"
	default:
		return "<XMPPMessage" +
			" id=" + this.id +
			" type=" + fmt.Sprint(this.Type()) +
			" value=" + fmt.Sprint(this.Value()) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// EQUALS

func (this *XMPPMessage) Equals(other *XMPPMessage) bool {
	// Never matches if types are different
	if this.Type() != other.Type() {
		return false
	}
	// Try and match values which are either uint, string or array
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
	case []interface{}:
		thisValues := this.Value().([]interface{})
		otherValues := other.Value().([]interface{})
		if len(thisValues) != len(otherValues) {
			return false
		}
		for i := 0; i < len(thisValues); i++ {
			if thisValues[i] != otherValues[i] {
				return false
			}
		}
		return true
	}

	// In default case, indicate no equals
	return false
}
