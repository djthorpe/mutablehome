/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	mutablehome "github.com/djthorpe/mutablehome"
	coap "github.com/go-ocf/go-coap"
	codes "github.com/go-ocf/go-coap/codes"
	dtls "github.com/pion/dtls/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tradfri struct {
	Id      string
	Key     string
	Path    string
	Timeout time.Duration
}

type tradfri struct {
	log     gopi.Logger
	id      string
	key     string
	path    string
	timeout time.Duration
	conn    *coap.ClientConn
	token   *token
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	TIMEOUT            = 5 * time.Second
	PATH_AUTH_EXCHANGE = "/15011/9063"
	PATH_DEVICES       = "/15001"
	PATH_GROUPS        = "/15004"
	PATH_SCENES        = "/15005"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Tradfri) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<mutablehome.ikea-tradfri.Open>{ config=%+v }", config)

	this := new(tradfri)
	this.log = logger
	this.id = config.Id
	this.key = config.Key
	this.token = &token{}

	// Set timeout
	if config.Timeout == 0 {
		this.timeout = TIMEOUT
	} else {
		this.timeout = config.Timeout
	}

	// Create path if it doesn't exist, and read token file
	if path, err := this.createPath(config.Path); err != nil {
		return nil, err
	} else if err := this.token.Read(path); err != nil {
		return nil, err
	} else {
		this.path = path
	}

	// Success
	return this, nil
}

func (this *tradfri) Close() error {
	this.log.Debug("<mutablehome.ikea-tradfri.Close>{}")

	// Close connection
	if this.conn != nil {
		if err := this.conn.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.conn = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *tradfri) String() string {
	if this.conn == nil {
		return fmt.Sprintf("<mutablehome.ikea-tradfri>{ nil }")
	} else {
		return fmt.Sprintf("<mutablehome.ikea-tradfri>{ addr=%v }", this.conn.RemoteAddr())
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *tradfri) Connect(service gopi.RPCServiceRecord, flags gopi.RPCFlag) error {
	this.log.Debug2("<mutablehome.ikea-tradfri.Connect>{ service=%v }", service)

	if addr, err := this.getAddr(service, flags); err != nil {
		return err
	} else {
		if this.token.Id == "" || this.token.Token == "" {
			if id, err := this.getId(service); err != nil {
				return err
			} else if conn, err := this.connectWith(addr, "Client_identity", this.key); err != nil {
				return err
			} else if token, err := this.authenticate(conn, id); err != nil {
				return err
			} else if err := token.Write(this.path); err != nil {
				return err
			} else if err := conn.Close(); err != nil {
				return err
			}
		}

		// Connect with existing token parameters
		if conn, err := this.connectWith(addr, this.token.Id, this.token.Token); err != nil {
			return err
		} else {
			this.conn = conn
		}
	}

	// Success
	return nil
}

func (this *tradfri) connectWith(addr, name, key string) (*coap.ClientConn, error) {
	if conn, err := coap.DialDTLSWithTimeout("udp", addr, &dtls.Config{
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

func (this *tradfri) authenticate(conn *coap.ClientConn, id string) (*token, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Authenticate>{ id=%v }", strconv.Quote(id))

	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	body := strings.NewReader(fmt.Sprintf(`{"9090":"%s"}`, id))
	if response, err := conn.PostWithContext(ctx, PATH_AUTH_EXCHANGE, coap.AppJSON, body); err != nil {
		return nil, err
	} else if response.Code() != codes.Created {
		return nil, fmt.Errorf("%w: %v", gopi.ErrUnexpectedResponse, response.Code())
	} else {
		token := &token{}
		payload := response.Payload()
		payload = payload[0 : len(payload)-3]
		if err := json.Unmarshal(payload, &token); err != nil {
			return nil, fmt.Errorf("%w: %v", err, string(payload))
		} else {
			token.Id = id
			return token, nil
		}
	}
}

func (this *tradfri) Devices() ([]uint, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Devices>{ }")
	return this.requestIdsForPath(PATH_DEVICES)
}

func (this *tradfri) Groups() ([]uint, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Groups>{ }")
	return this.requestIdsForPath(PATH_GROUPS)
}

func (this *tradfri) Scenes() ([]uint, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Scenes>{ }")
	return this.requestIdsForPath(PATH_SCENES)
}

func (this *tradfri) Device(id uint) (mutablehome.IkeaDevice, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Device>{ id=%v }", id)

	device := &device{}
	if err := this.requestObjForPathId(PATH_DEVICES, id, device); err != nil {
		return nil, err
	} else {
		return device, nil
	}
}

func (this *tradfri) Group(id uint) (mutablehome.IkeaGroup, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Group>{ id=%v }", id)

	group := &group{}
	if err := this.requestObjForPathId(PATH_GROUPS, id, group); err != nil {
		return nil, err
	} else {
		return group, nil
	}
}

func (this *tradfri) Scene(id uint) (mutablehome.IkeaScene, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.Scene>{ id=%v }", id)

	scene := &scene{}
	if err := this.requestObjForPathId(PATH_SCENES, id, scene); err != nil {
		return nil, err
	} else {
		return scene, nil
	}
}

func (this *tradfri) ObserveDevice(ctx context.Context, id uint) error {
	path := fmt.Sprintf("%v/%d", PATH_DEVICES, id)
	if obs, err := this.conn.Observe(path, func(response *coap.Request) {
		device := &device{}
		if response.Msg.Code() != codes.Content {
			fmt.Println(fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Msg.Code(), response.Msg.Path()))
		} else if err := json.Unmarshal(response.Msg.Payload(), &device); err != nil {
			fmt.Println(fmt.Errorf("%w: %v", err, string(response.Msg.Payload())))
		} else {
			fmt.Println("Got", device)
		}
	}); err != nil {
		return err
	} else {
		select {
		case <-ctx.Done():
			if err := obs.Cancel(); err != nil {
				return err
			}
		}
	}
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *tradfri) createPath(path string) (string, error) {
	// If path is relative, then append user's home folder
	if filepath.IsAbs(path) == false {
		if home, err := os.UserHomeDir(); err != nil {
			return "", err
		} else {
			path = filepath.Join(home, path)
		}
	}
	// If path doesn't exist then try and create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			return path, err
		}
	}
	// Make sure path is available
	if stat, err := os.Stat(path); err != nil {
		return path, err
	} else if stat.IsDir() == false {
		return path, fmt.Errorf("%w: Not a folder: %v", gopi.ErrBadParameter, path)
	}
	// Success
	return path, nil
}

func (this *tradfri) getId(service gopi.RPCServiceRecord) (string, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.SetId>{ service=%v }", service)

	id := this.id
	if service == nil {
		return "", gopi.ErrBadParameter
	} else if id == "" {
		id = service.Name()
	}

	if id = strings.TrimSpace(id); len(id) == 0 {
		return "", fmt.Errorf("%w: Missing or invalid ID", gopi.ErrBadParameter)
	} else {
		return id, nil
	}
}

func (this *tradfri) getAddr(service gopi.RPCServiceRecord, flag gopi.RPCFlag) (string, error) {
	this.log.Debug2("<mutablehome.ikea-tradfri.getAddr>{ service=%v flag=%v }", service, flag)

	switch flag & (gopi.RPC_FLAG_INET_V4 | gopi.RPC_FLAG_INET_V6) {
	case gopi.RPC_FLAG_INET_V4:
		if addr := service.IP4(); len(addr) > 0 && service.Port() > 0 {
			return fmt.Sprintf("%s:%d", addr[0].String(), service.Port()), nil
		}
	case gopi.RPC_FLAG_INET_V6:
		if addr := service.IP6(); len(addr) > 0 && service.Port() > 0 {
			return fmt.Sprintf("[%s]:%d", addr[0].String(), service.Port()), nil
		}
	case gopi.RPC_FLAG_INET_V4 | gopi.RPC_FLAG_INET_V6:
		if addr := service.IP4(); len(addr) > 0 && service.Port() > 0 {
			return fmt.Sprintf("%s:%d", addr[0].String(), service.Port()), nil
		}
		if addr := service.IP6(); len(addr) > 0 && service.Port() > 0 {
			return fmt.Sprintf("[%s]:%d", addr[0].String(), service.Port()), nil
		}
	}

	// Success
	return "", fmt.Errorf("%w: Cannot determine gateway address and/or port", gopi.ErrBadParameter)
}

func (this *tradfri) requestIdsForPath(path string) ([]uint, error) {

	if this.conn == nil {
		return nil, fmt.Errorf("%w: Not connected", gopi.ErrNotFound)
	}

	var ids []uint
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()

	if response, err := this.conn.GetWithContext(ctx, path); err != nil {
		return nil, err
	} else if response.Code() != codes.Content {
		return nil, fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), path)
	} else {
		if err := json.Unmarshal(response.Payload(), &ids); err != nil {
			return nil, fmt.Errorf("%w: %v", err, string(response.Payload()))
		} else {
			return ids, nil
		}
	}
}

func (this *tradfri) requestObjForPathId(path string, id uint, obj interface{}) error {
	if this.conn == nil {
		return fmt.Errorf("%w: Not connected", gopi.ErrNotFound)
	}
	ctx, cancel := context.WithTimeout(context.Background(), this.timeout)
	defer cancel()
	if response, err := this.conn.GetWithContext(ctx, fmt.Sprintf("%v/%d", path, id)); err != nil {
		return err
	} else if response.Code() != codes.Content {
		return fmt.Errorf("%w: %v (path: %v)", gopi.ErrUnexpectedResponse, response.Code(), response.Path())
	} else if err := json.Unmarshal(response.Payload(), obj); err != nil {
		return fmt.Errorf("%w: %v", err, string(response.Payload()))
	}
	// Success
	return nil
}
