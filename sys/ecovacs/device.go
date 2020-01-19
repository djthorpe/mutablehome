package ecovacs

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
	xmpp "github.com/mattn/go-xmpp"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	DeviceId  string `json:"did"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	Resource  string `json:"resource"`
	Nickname_ string `json:"nick"`
	Company   string `json:"company"`

	source *ecovacs
	client *xmpp.Client

	reqId   uint
	reqLock sync.Mutex
	stop    chan struct{}
	state

	sync.Mutex
	sync.WaitGroup
}

type message struct {
	XMLName xml.Name `xml:"query"`
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
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Time-to-live for any state
	DELTA_TTL = 4 * time.Minute
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	ServerMap = map[string]string{
		"msg.ecouser.net":    "CH",
		"msg-as.ecouser.net": "TW,MY,JP,SG,TH,HK,IN,KR",
		"msg-na.ecouser.net": "US",
		"msg-eu.ecouser.net": "FR,ES,UK,NO,MX,DE,PT,CH,AU,IT,NL,SE,BE,DK",
		"msg-ww.ecouser.net": "",
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *device) Connect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.client != nil {
		return gopi.ErrInternalAppError.WithPrefix("Connect")
	}
	if server := CountryToXMPPServer(this.source.country); server == "" {
		return gopi.ErrInternalAppError.WithPrefix("country")
	} else {
		opts := xmpp.Options{
			Host:     fmt.Sprintf("%s:%d", server, ECOVACS_XMPP_PORT),
			User:     fmt.Sprintf("%s@%s", this.source.userId, ECOVACS_REALM),
			Password: fmt.Sprintf("0/%s/%s", this.source.resourceId, this.source.accessToken),
			NoTLS:    true,
			Session:  true,

			Debug: this.source.Log.IsDebug(),
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		if client, err := opts.NewClient(); err != nil {
			return err
		} else {
			this.client = client
			this.stop = make(chan struct{})
			this.WaitGroup.Add(2)
			go this.recv()
			go this.ping(this.stop)
		}
	}
	// Success
	return nil
}

func (this *device) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// If client is nil then no need to disconnect
	if this.client == nil {
		return nil
	}

	// end ping
	close(this.stop)

	// close client
	err := this.client.Close()

	// wait for termination of recv
	this.WaitGroup.Wait()

	// release resources
	this.client = nil
	this.stop = nil

	// return any error condition
	return err
}

////////////////////////////////////////////////////////////////////////////////
// SEND COMMANDS

func (this *device) Ping() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.client.PingC2S(this.Address(), this.client.JID())
}

func (this *device) Clean(mode home.EcovacsCleanMode, suction home.EcovacsCleanSuction) (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := fmt.Sprintf(`<ctl td="Clean"><clean type="%s" speed="%s"/></ctl>`, string(mode), string(suction))
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("Clean")
	} else {
		return this.send(command)
	}
}

func (this *device) Charge() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="Charge"><charge type="go"/></ctl>`
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("Charge")
	} else {
		return this.send(command)
	}
}

func (this *device) GetBatteryInfo() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetBatteryInfo"></ctl>`
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetBatteryInfo")
	} else {
		return this.send(command)
	}
}

func (this *device) GetLifeSpan(part home.EcovacsPart) (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := fmt.Sprintf(`<ctl td="GetLifeSpan" type="%s"></ctl>`, string(part))
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetLifeSpan")
	} else {
		return this.send(command)
	}
}

func (this *device) GetChargeState() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetChargeState"></ctl>`
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetChargeState")
	} else {
		return this.send(command)
	}
}

func (this *device) GetCleanState() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetCleanState"></ctl>`
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetCleanState")
	} else {
		return this.send(command)
	}
}

func (this *device) GetVersion() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetVersion" name="FW"></ctl>`
	if this.client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetVersion")
	} else {
		return this.send(command)
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *device) Address() string {
	return fmt.Sprintf("%s@%s.ecorobot.net/atom", this.DeviceId, this.Class)
}

func (this *device) Nickname() string {
	return this.Nickname_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return "<ecovacs.device" +
		" addr=" + strconv.Quote(this.Address()) +
		" nickname=" + strconv.Quote(this.Nickname_) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *device) nextRequestId() string {
	this.reqLock.Lock()
	defer this.reqLock.Unlock()
	this.reqId = this.reqId + 1
	return fmt.Sprint(this.reqId)
}

func (this *device) send(command string) (string, error) {
	return this.client.RawInformationQuery(this.client.JID(), this.Address(), this.nextRequestId(), xmpp.IQTypeSet, "com:ctl", command)
}

func (this *device) recv() {
	defer this.WaitGroup.Done()
FOR_LOOP:
	for {
		if stanza, err := this.client.Recv(); err != nil {
			want := "use of closed network connection"
			if strings.Contains(err.Error(), want) {
				break FOR_LOOP
			} else {
				fmt.Println("RECV ERROR", err)
			}
		} else {
			switch v := stanza.(type) {
			case xmpp.IQ:
				// If message is from this device, then handle
				if v.From == this.Address() {
					if err := this.recv_parse(v.ID, v.Type, v.Query); err != nil {
						fmt.Println("PARSE ERROR", err, string(v.Query))
					}
				}
			}
		}
	}
}

func (this *device) recv_parse(reqId string, type_ string, data []byte) error {
	var m message
	if len(data) > 0 {
		if err := xml.Unmarshal(data, &m); err != nil {
			return err
		}
	}
	if len(data) > 0 && type_ == "set" {
		evt := NewEvent(this.source, this, reqId, &m, data)
		if evt.Type() == home.ECOVACS_EVENT_NONE {
			// Emit but don't store NONE events
			this.source.bus.Emit(evt)
		} else if modified := this.state.Set(evt, DELTA_TTL); modified {
			// Emit but only if modified
			this.source.bus.Emit(evt)
			// Print out the current device state
			fmt.Println(this.state.String())
		}
	}
	return nil
}

func (this *device) updateStatusForKey(key home.EcovacsEventType) error {
	switch key {
	case home.ECOVACS_EVENT_BATTERYLEVEL:
		if _, err := this.GetBatteryInfo(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_CHARGESTATE:
		if _, err := this.GetChargeState(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_CLEANSTATE:
		if _, err := this.GetCleanState(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_LIFESPAN:
		if _, err := this.GetLifeSpan(home.ECOVACS_PART_BRUSH); err != nil {
			return err
		}
		if _, err := this.GetLifeSpan(home.ECOVACS_PART_DUSTFILTER); err != nil {
			return err
		}
		if _, err := this.GetLifeSpan(home.ECOVACS_PART_SIDEBRUSH); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_VERSION:
		if _, err := this.GetVersion(); err != nil {
			return err
		}
	default:
		return gopi.ErrBadParameter.WithPrefix(fmt.Sprint(key))
	}

	// Success
	return nil
}

func (this *device) ping(stop <-chan struct{}) {
	defer this.WaitGroup.Done()
	ping_ticker := time.NewTicker(time.Second * 30)
	update_ticker := time.NewTicker(time.Second * 15)
FOR_LOOP:
	for {
		select {
		case <-ping_ticker.C:
			if err := this.Ping(); err != nil {
				fmt.Println("PING ERROR", err)
			}
		case <-update_ticker.C:
			if key := this.state.NextExpiredKey(); key != home.ECOVACS_EVENT_NONE {
				if err := this.updateStatusForKey(key); err != nil {
					fmt.Println("UPDATE ERROR", err)
				}
			}
		case <-stop:
			ping_ticker.Stop()
			update_ticker.Stop()
			break FOR_LOOP
		}
	}
}

func CountryToXMPPServer(country string) string {
	d := ""
	for key, value := range ServerMap {
		if value == "" {
			d = key
		}
		for _, code := range strings.Split(value, ",") {
			if strings.ToUpper(country) == code {
				return key
			}
		}
	}
	return d
}
