package ecovacs

import (
	"fmt"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
	xmpp "github.com/mattn/go-xmpp"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type XMPPClient struct {
	Client   *xmpp.Client
	DeviceId string
	Class    string

	RequestId
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// NEW CLIENT / CLOSE

func (this *XMPPClient) NewClient(opts xmpp.Options, deviceId, class string) error {
	if client, err := opts.NewClient(); err != nil {
		return err
	} else {
		this.Client = client
		this.DeviceId = deviceId
		this.Class = class
	}

	// Return success
	return nil
}

func (this *XMPPClient) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.Client == nil {
		return nil
	}
	err := this.Client.Close()
	this.Client = nil
	return err
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *XMPPClient) IsConnected() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.Client != nil
}

func (this *XMPPClient) Address() string {
	return fmt.Sprintf("%s@%s.ecorobot.net/atom", this.DeviceId, this.Class)
}

////////////////////////////////////////////////////////////////////////////////
// SEND COMMANDS

func (this *XMPPClient) Ping() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.Client.PingC2S(this.Address(), this.Client.JID())
}

func (this *XMPPClient) Clean(mode home.EcovacsCleanMode, suction home.EcovacsCleanSuction) (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := fmt.Sprintf(`<ctl td="Clean"><clean type="%s" speed="%s"/></ctl>`, string(mode), string(suction))
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("Clean")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) Charge() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="Charge"><charge type="go"/></ctl>`
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("Charge")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) GetBatteryInfo() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetBatteryInfo"></ctl>`
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetBatteryInfo")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) GetLifeSpan(part home.EcovacsPart) (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := fmt.Sprintf(`<ctl td="GetLifeSpan" type="%s"></ctl>`, string(part))
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetLifeSpan")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) GetChargeState() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetChargeState"></ctl>`
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetChargeState")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) GetCleanState() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetCleanState"></ctl>`
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetCleanState")
	} else {
		return this.send(command)
	}
}

func (this *XMPPClient) GetVersion() (string, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	command := `<ctl td="GetVersion" name="FW"></ctl>`
	if this.Client == nil {
		return "", gopi.ErrInternalAppError.WithPrefix("GetVersion")
	} else {
		return this.send(command)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *XMPPClient) send(command string) (string, error) {
	return this.Client.RawInformationQuery(this.Client.JID(), this.Address(), this.RequestId.Next(), xmpp.IQTypeSet, "com:ctl", command)
}

func (this *XMPPClient) Recv() (*XMPPMessage, error) {
	for {
		// Receive a message from XMPP
		stanza, err := this.Client.Recv()

		// Deal with "closed network connection" error, which is
		// due to us closing the intenet connection and we don't need
		// to report the error but simply return nil. Other error here
		// is probably severe and we should disconnect from the client
		if err != nil {
			want := "use of closed network connection"
			if strings.Contains(err.Error(), want) {
				// Error due to disconnect, don't report
				return nil, nil
			} else {
				// Other error should result in disconnect/connect cycle
				return nil, err
			}
		}
		switch v := stanza.(type) {
		case xmpp.IQ:
			// If message is from this device, then handle
			if message, err := Parse(v); err != nil {
				fmt.Println("PARSE ERROR", err, string(v.Query))
			} else if message != nil {
				return message, nil
			}
		}
	}
}

func Parse(in xmpp.IQ) (*XMPPMessage, error) {
	if len(in.Query) == 0 || in.Type != "set" {
		return nil, nil
	} else {
		return NewXMPPMessage(in.Query, in.ID)
	}
}
