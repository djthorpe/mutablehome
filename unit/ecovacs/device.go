package ecovacs

import (
	"crypto/tls"
	"fmt"
	"strconv"
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
	DeviceId_ string `json:"did"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	Resource  string `json:"resource"`
	Nickname_ string `json:"nick"`
	Company   string `json:"company"`

	source *ecovacs
	stop   chan struct{}

	XMPPClient
	DeviceState
	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Time-to-live for any state
	DELTA_OTHER_TTL    = 4 * time.Minute
	DELTA_VERSION_TTL  = 6 * time.Hour
	DELTA_LIFESPAN_TTL = 2 * time.Hour
)

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *device) Connect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.XMPPClient.IsConnected() {
		return gopi.ErrInternalAppError.WithPrefix("Connect")
	}
	if server := CountryToXMPPServer(this.source.country); server == "" {
		return gopi.ErrInternalAppError.WithPrefix("country")
	} else if err := this.XMPPClient.NewClient(xmpp.Options{
		Host:     fmt.Sprintf("%s:%d", server, ECOVACS_XMPP_PORT),
		User:     fmt.Sprintf("%s@%s", this.source.userId, ECOVACS_REALM),
		Password: fmt.Sprintf("0/%s/%s", this.source.resourceId, this.source.accessToken),
		NoTLS:    true,
		Session:  true,

		Debug: this.source.Log.IsDebug(),
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}, this.DeviceId_, this.Class); err != nil {
		return err
	} else {
		this.stop = make(chan struct{})
		this.WaitGroup.Add(2)
		go this.recv()
		go this.ping(this.stop)
	}
	// Return success
	return nil
}

func (this *device) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// If client is nil then no need to disconnect
	if this.XMPPClient.IsConnected() == false {
		return nil
	}

	// end ping
	close(this.stop)

	// close client
	err := this.XMPPClient.Close()

	// wait for termination of recv
	this.WaitGroup.Wait()

	// release resources
	this.stop = nil

	// return any error condition
	return err
}

////////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *device) Id() string {
	return this.DeviceId_
}

func (this *device) Nickname() string {
	return this.Nickname_
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return "<ecovacs.device" +
		" device_id=" + strconv.Quote(this.DeviceId_) +
		" nickname=" + strconv.Quote(this.Nickname_) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *device) recv() {
	defer this.WaitGroup.Done()
FOR_LOOP:
	for {
		if message, err := this.XMPPClient.Recv(); err != nil {
			// We need to do this in a goroutine to prevent deadlock
			go func() {
				this.source.deviceError(this, err)
			}()
			break FOR_LOOP
		} else if message == nil {
			// End of cycle when no message returned
			break FOR_LOOP
		} else {
			// Create event from message
			event := NewEvent(this.source, this, message)
			if event.Type() == home.ECOVACS_EVENT_NONE {
				// Emit but don't store messages which can't be decoded
				this.source.bus.Emit(event)
			} else if modified := this.DeviceState.Set(message, ttlForType(event.Type())); modified {
				// Emit messages but only if modified
				this.source.bus.Emit(event)
			}
		}
	}
}

func (this *device) updateStatusForKey(key home.EcovacsEventType) error {
	switch key {
	case home.ECOVACS_EVENT_BATTERYLEVEL:
		if _, err := this.XMPPClient.GetBatteryInfo(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_CHARGESTATE:
		if _, err := this.XMPPClient.GetChargeState(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_CLEANSTATE:
		if _, err := this.XMPPClient.GetCleanState(); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_LIFESPAN:
		if _, err := this.XMPPClient.GetLifeSpan(home.ECOVACS_PART_BRUSH); err != nil {
			return err
		}
		if _, err := this.XMPPClient.GetLifeSpan(home.ECOVACS_PART_DUSTFILTER); err != nil {
			return err
		}
		if _, err := this.XMPPClient.GetLifeSpan(home.ECOVACS_PART_SIDEBRUSH); err != nil {
			return err
		}
	case home.ECOVACS_EVENT_VERSION:
		if _, err := this.XMPPClient.GetVersion(); err != nil {
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

	// Ping every 30 seconds and update device state every 15 seconds
	ping_ticker := time.NewTicker(time.Second * 30)
	update_ticker := time.NewTicker(time.Second * 15)

	// Add expired keys to ticker
	this.DeviceState.AddExpiredKey(home.ECOVACS_EVENT_BATTERYLEVEL)
	this.DeviceState.AddExpiredKey(home.ECOVACS_EVENT_CHARGESTATE)
	this.DeviceState.AddExpiredKey(home.ECOVACS_EVENT_CLEANSTATE)
	this.DeviceState.AddExpiredKey(home.ECOVACS_EVENT_LIFESPAN)
	this.DeviceState.AddExpiredKey(home.ECOVACS_EVENT_VERSION)

FOR_LOOP:
	for {
		select {
		case <-ping_ticker.C:
			if err := this.XMPPClient.Ping(); err != nil {
				fmt.Println("PING ERROR", err)
			}
		case <-update_ticker.C:
			if key := this.DeviceState.NextExpiredKey(); key != home.ECOVACS_EVENT_NONE {
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

func ttlForType(type_ home.EcovacsEventType) time.Duration {
	switch type_ {
	case home.ECOVACS_EVENT_VERSION:
		return DELTA_VERSION_TTL
	case home.ECOVACS_EVENT_LIFESPAN:
		return DELTA_LIFESPAN_TTL
	default:
		return DELTA_OTHER_TTL
	}
}
