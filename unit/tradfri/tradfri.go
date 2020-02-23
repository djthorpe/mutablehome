/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	"github.com/djthorpe/mutablehome"
	"github.com/go-ocf/go-coap"
	"github.com/go-ocf/go-coap/codes"
	"github.com/pion/dtls/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tradfri struct {
	Id      string
	Key     string
	Path    string
	Timeout time.Duration
	Bus     gopi.Bus
}

type tradfri struct {
	log     gopi.Logger
	bus     gopi.Bus
	id      string
	key     string
	path    string
	timeout time.Duration
	conn    *coap.ClientConn
	devices map[uint]*device

	Token
	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	CONN_TIMEOUT       = 5 * time.Second
	PATH_AUTH_EXCHANGE = "/15011/9063"
	PATH_DEVICES       = "/15001"
	PATH_GROUPS        = "/15004"
	PATH_SCENES        = "/15005"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Tradfri) Name() string { return "tradfri" }

func (config Tradfri) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(tradfri)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *tradfri) Init(config Tradfri) error {
	this.id = config.Id
	this.key = config.Key

	// Check for bus
	if config.Bus == nil {
		return gopi.ErrBadParameter.WithPrefix("bus")
	} else {
		this.bus = config.Bus
	}

	// Set timeout
	if config.Timeout == 0 {
		this.timeout = CONN_TIMEOUT
	} else {
		this.timeout = config.Timeout
	}

	// Create path if it doesn't exist, and read token file
	if path, err := this.Token.CreatePath(config.Path); err != nil {
		return err
	} else if err := this.Token.Read(path); err != nil {
		return err
	} else {
		this.path = path
	}

	// Create empty devices map
	this.devices = make(map[uint]*device)

	// Success
	return nil
}

func (this *tradfri) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close connection
	if this.conn != nil {
		if err := this.conn.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.conn = nil
	this.devices = nil
	this.bus = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *tradfri) Connect(service gopi.RPCServiceRecord, flags gopi.RPCFlag) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.conn != nil {
		return gopi.ErrOutOfOrder
	} else if addr, err := this.getAddr(service, flags); err != nil {
		return err
	} else {
		// Authenticate and close connection
		if this.Token.Id == "" || this.Token.Token == "" {
			if id, err := this.getId(service); err != nil {
				return err
			} else if conn, err := this.connectWith(addr, "Client_identity", this.key); err != nil {
				return err
			} else if err := this.authenticate(conn, id); err != nil {
				return err
			} else if err := this.Token.Write(this.path); err != nil {
				return err
			} else if err := conn.Close(); err != nil {
				return err
			}
		}

		// Connect with existing token parameters
		if conn, err := this.connectWith(addr, this.Token.Id, this.Token.Token); err != nil {
			return err
		} else {
			this.conn = conn
		}
	}

	// Success
	return nil
}

func (this *tradfri) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close connection
	if this.conn != nil {
		if err := this.conn.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.conn = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *tradfri) String() string {
	return "<" + this.Log.Name() +
		" id=" + strconv.Quote(this.id) +
		" key=" + strconv.Quote(this.key) +
		" path=" + strconv.Quote(this.path) +
		" timeout=" + fmt.Sprint(this.timeout) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// DEVICES, GROUPS AND SCENES

func (this *tradfri) Devices() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_DEVICES)
}

func (this *tradfri) Groups() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_GROUPS)
}

func (this *tradfri) Scenes() ([]uint, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return this.requestIdsForPath(PATH_SCENES)
}

func (this *tradfri) Device(id uint) (mutablehome.IkeaDevice, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	device := NewDevice()
	if err := this.requestObjForPathId(PATH_DEVICES, id, device); err != nil {
		return nil, err
	} else {
		return device, nil
	}
}

func (this *tradfri) Group(id uint) (mutablehome.IkeaGroup, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	group := NewGroup()
	if err := this.requestObjForPathId(PATH_GROUPS, id, group); err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SEND COMMANDS TO GATEWAY

func (this *tradfri) Send(commands ...mutablehome.IkeaCommand) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(commands) == 0 {
		return gopi.ErrBadParameter.WithPrefix("values")
	}
	for _, command := range commands {
		if command == nil {
			return gopi.ErrBadParameter.WithPrefix("command")
		} else if body, err := command.Body(); err != nil {
			return gopi.ErrBadParameter.WithPrefix("command")
		} else if message, err := this.conn.Put(command.Path(), coap.AppJSON, body); err != nil {
			return err
		} else if message.Code() != codes.Changed {
			return fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, message.Code(), message.Path())
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// OBSERVE DEVICE CHANGES

func (this *tradfri) ObserveDevice(ctx context.Context, id uint) error {
	path := strings.Join([]string{PATH_DEVICES, fmt.Sprint(id)}, "/")
	ticker := time.NewTimer(100 * time.Millisecond)
	var obs *coap.Observation

FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			// Stop
			if obs != nil {
				if err := obs.Cancel(); err != nil {
					return err
				}
			}
			// Start
			if obs_, err := this.conn.Observe(path, this.observeDeviceCallback); err != nil {
				return err
			} else {
				obs = obs_
			}
			// Restart the ticker with some random additional interval
			ticker.Reset(time.Second * (5 + time.Duration(rand.Int31n(15))))
		case <-ctx.Done():
			if err := obs.Cancel(); err != nil {
				return err
			} else {
				break FOR_LOOP
			}
		}
	}

	// Success
	return ctx.Err()
}

func (this *tradfri) observeDeviceCallback(response *coap.Request) {
	device := NewDevice()
	if response.Msg.Code() != codes.Content {
		this.Log.Error(fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Msg.Code()))
		return
	} else if err := json.Unmarshal(response.Msg.Payload(), &device); err != nil {
		this.Log.Error(fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, err))
		return
	}

	// Emit device event if there is a change
	if event := this.observeDeviceEvent(device); event != mutablehome.IKEA_EVENT_NONE {
		this.bus.Emit(NewEvent(this, event, device))
	}
}

func (this *tradfri) observeDeviceEvent(device *device) mutablehome.IkeaEventType {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device == nil || device.Id() == 0 {
		return mutablehome.IKEA_EVENT_NONE
	}

	// Check for added devices
	id := device.Id()
	if other, exists := this.devices[id]; exists == false {
		this.devices[id] = device
		return mutablehome.IKEA_EVENT_DEVICE_ADDED
	} else if device.Equals(other) == false {
		this.devices[id] = device
		return mutablehome.IKEA_EVENT_DEVICE_CHANGED
	} else {
		return mutablehome.IKEA_EVENT_NONE
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// ConnectWith creates a COAP connection
func (this *tradfri) connectWith(addr, name, key string) (*coap.ClientConn, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: Missing key parameter", gopi.ErrBadParameter)
	} else if conn, err := coap.DialDTLSWithTimeout("udp", addr, &dtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			return []byte(key), nil
		},
		PSKIdentityHint: []byte(name),
		CipherSuites:    []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	}, this.timeout); err != nil {
		return nil, fmt.Errorf("%w (addr: %s)", err, addr)
	} else {
		return conn, nil
	}
}

// getAddr returns an address to connect to from a service record
func (this *tradfri) getAddr(service gopi.RPCServiceRecord, flag gopi.RPCFlag) (string, error) {
	for _, addr := range service.Addrs {
		switch {
		case flag&gopi.RPC_FLAG_INET_V4 == gopi.RPC_FLAG_INET_V4:
			if addr.To4() != nil && service.Port > 0 {
				return fmt.Sprintf("%s:%d", addr.String(), service.Port), nil
			}
		case flag&gopi.RPC_FLAG_INET_V6 == gopi.RPC_FLAG_INET_V6:
			if addr.To16() != nil && service.Port > 0 {
				return fmt.Sprintf("[%s]:%d", addr.String(), service.Port), nil
			}
		}
	}

	// Return error
	return "", fmt.Errorf("%w: Cannot determine gateway address and/or port", gopi.ErrBadParameter)
}

// getId returns a unique id
func (this *tradfri) getId(service gopi.RPCServiceRecord) (string, error) {
	id := this.id
	if id == "" {
		id = service.Name
	}

	if id = strings.TrimSpace(id); len(id) == 0 {
		return "", fmt.Errorf("%w: Missing or invalid ID", gopi.ErrBadParameter)
	} else {
		return id, nil
	}
}

// authenticate performs the gateway authentication
func (this *tradfri) authenticate(conn *coap.ClientConn, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	body := strings.NewReader(fmt.Sprintf(`{"9090":"%s"}`, id))
	if response, err := conn.PostWithContext(ctx, PATH_AUTH_EXCHANGE, coap.AppJSON, body); err != nil {
		return err
	} else if response.Code() != codes.Created {
		return fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Code())
	} else {
		payload := response.Payload()
		payload = payload[0 : len(payload)-3]
		if err := json.Unmarshal(payload, &this.Token); err != nil {
			return fmt.Errorf("%w: %v", err, string(payload))
		}
	}

	// Set token id
	this.Token.Id = id

	// Return success
	return nil
}

func (this *tradfri) requestIdsForPath(path string) ([]uint, error) {
	if this.conn == nil {
		return nil, gopi.ErrOutOfOrder
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	var ids []uint
	if response, err := this.conn.GetWithContext(ctx, path); err != nil {
		return nil, err
	} else if response.Code() != codes.Content {
		return nil, fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), path)
	} else if err := json.Unmarshal(response.Payload(), &ids); err != nil {
		return nil, fmt.Errorf("%w: %v", err, string(response.Payload()))
	}

	// Success
	return ids, nil
}

func (this *tradfri) requestObjForPathId(path string, id uint, obj interface{}) error {
	if this.conn == nil {
		return gopi.ErrOutOfOrder
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()
	if response, err := this.conn.GetWithContext(ctx, fmt.Sprintf("%v/%d", path, id)); err != nil {
		return err
	} else if response.Code() != codes.Content {
		return fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), response.Path())
	} else if err := json.Unmarshal(response.Payload(), obj); err != nil {
		return fmt.Errorf("%w: %v", err, string(response.Payload()))
	} else {
		fmt.Println(string(response.Payload()))
	}

	// Success
	return nil
}
