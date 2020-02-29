/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	iface "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type devices struct {
	devices map[string]*device
	sync.Mutex
}

type Device struct {
	Service gopi.RPCServiceRecord
}

type device struct {
	service gopi.RPCServiceRecord
	txt     map[string]string
	stop    chan struct{}

	// chromecast state
	volume *volume
	app    *application

	connection
	channel
	sync.Mutex
	sync.WaitGroup
	base.Unit
}

const (
	READ_TIMEOUT    = 500 * time.Millisecond
	STATUS_INTERVAL = 10 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Device) Name() string { return "googlecast/device" }

func (config Device) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(device)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device

func (this *device) Init(config Device) error {
	// Set service
	this.service = config.Service
	this.channel.C = make(chan interface{}, 10)

	// Return success
	return nil
}

func (this *device) Close() error {
	// Disconnect
	if err := this.Disconnect(); err != nil {
		return err
	}

	// Close channel
	close(this.channel.C)

	// Release resources
	this.txt = nil
	this.channel.C = nil

	return this.Unit.Close()
}

func (this *device) Connect(flags gopi.RPCFlag, timeout time.Duration) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := this.connection.Connect(this.service, flags, timeout); err != nil {
		return err
	} else {
		this.stop = make(chan struct{})
		go this.rcv(this.stop)
	}

	if data, err := this.channel.Connect(); err != nil {
		return err
	} else {
		return this.send(data)
	}
}

func (this *device) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() {
		// Send disconnect to channel
		if data, err := this.channel.Disconnect(); err != nil {
			return err
		} else if err := this.send(data); err != nil {
			return err
		}

		// Stop receiving messages
		close(this.stop)
		this.WaitGroup.Wait()

		// Disconnect
		if err := this.connection.Disconnect(); err != nil {
			return err
		}
	}

	// Release resources
	this.stop = nil
	this.volume = nil
	this.app = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device VOLUME

func (this *device) Volume() iface.CastVolume {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() == false {
		return nil
	} else if this.volume == nil {
		return nil
	} else {
		return this.volume
	}
}

func (this *device) SetVolume(level float32) error {
	if level == 0 {
		return this.SetVolumeEx(0.0, true)
	} else {
		return this.SetVolumeEx(level, false)
	}
}

func (this *device) SetMute(mute bool) error {
	if this.volume == nil {
		return this.SetVolumeEx(0.5, mute)
	} else {
		return this.SetVolumeEx(this.volume.Level_, mute)
	}
}

func (this *device) SetVolumeEx(level float32, muted bool) error {
	v := volume{level, muted}
	if this.connection.IsConnected() == false {
		return gopi.ErrOutOfOrder
	} else if level < 0.0 || level > 1.0 {
		return gopi.ErrBadParameter.WithPrefix("level")
	} else if data, err := this.channel.SetVolume(v); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	} else {
		this.setStateVolume(v)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device APPLICATION

func (this *device) App() iface.CastApp {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() == false {
		return nil
	} else if this.app == nil || this.app.AppId == "" {
		return nil
	} else {
		return this.app
	}
}

func (this *device) LaunchAppWithId(appId string) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() == false {
		return gopi.ErrOutOfOrder
	} else if data, err := this.channel.LaunchAppWithId(appId); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device PLAY, PAUSE and STOP

func (this *device) SetPlay(state bool) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() == false {
		return gopi.ErrOutOfOrder
	} else if data, err := this.channel.PlayStop(state); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *device) SetPause(state bool) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.connection.IsConnected() == false {
		return gopi.ErrOutOfOrder
	} else if data, err := this.channel.PlayPause(state == false); err != nil {
		return err
	} else if err := this.send(data); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE MESSAGES

func (this *device) rcv(stop <-chan struct{}) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	statusTimer := time.NewTimer(500 * time.Millisecond)

FOR_LOOP:
	for {
		select {
		case state := <-this.channel.C:
			switch state.(type) {
			case volume:
				this.setStateVolume(state.(volume))
			case application:
				this.setStateApplication(state.(application))
			default:
				this.Log.Warn(this.Name()+":", "Unhandled state change: ", state)
			}
		case <-statusTimer.C:
			// Update status if volume is nil
			if this.volume == nil || this.app == nil {
				if data, err := this.channel.GetStatus(); err != nil {
					this.Log.Warn("GetStatus: %v", err)
				} else if err := this.send(data); err != nil {
					this.Log.Warn("GetStatus: %v", err)
				}
			}

			// Update receiver status if empty
			statusTimer.Reset(STATUS_INTERVAL)
		case <-stop:
			// Break loop and stop timers
			statusTimer.Stop()
			break FOR_LOOP
		default:
			var length uint32
			if err := this.connection.conn.SetReadDeadline(time.Now().Add(READ_TIMEOUT)); err != nil {
				this.Log.Error(err)
			} else if err := binary.Read(this.conn, binary.BigEndian, &length); err != nil {
				if err == io.EOF || os.IsTimeout(err) {
					// Ignore error
				} else {
					this.Log.Error(err)
				}
			} else if length == 0 {
				this.Log.Error(fmt.Errorf("Received zero-sized data"))
			} else {
				payload := make([]byte, length)
				if bytes_read, err := io.ReadFull(this.conn, payload); err != nil {
					this.Log.Error(err)
				} else if bytes_read != int(length) {
					this.Log.Error(fmt.Errorf("Received different number of bytes %v read, expected %v", bytes_read, length))
				} else if data, err := this.channel.decode(payload); err != nil {
					this.Log.Error(err)
				} else if err := this.send(data); err != nil {
					this.Log.Error(err)
				}
			}
		}
	}
}

func (this *device) send(data []byte) error {
	if len(data) == 0 {
		return nil
	} else if err := binary.Write(this.connection.conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	} else if _, err := this.conn.Write(data); err != nil {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *device) Id() string {
	return this.Txt("id")
}

func (this *device) Name() string {
	return this.Txt("fn")
}

func (this *device) Model() string {
	return this.Txt("md")
}

func (this *device) Service() string {
	return this.Txt("rs")
}

func (this *device) State() uint {
	if value := this.Txt("st"); value == "" {
		return 0
	} else if value_, err := strconv.ParseUint(value, 10, 32); err != nil {
		return 0
	} else {
		return uint(value_)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	str := "<cast.Device" +
		" id=" + strconv.Quote(this.Id()) +
		" name=" + strconv.Quote(this.Name()) +
		" model=" + strconv.Quote(this.Model()) +
		" state=" + fmt.Sprint(this.State())

	if this.Service() != "" {
		str += " service=" + strconv.Quote(this.Service())
	}
	if this.connection.IsConnected() {
		str += " conn=" + this.connection.String()
	}
	if this.volume != nil {
		str += " volume=" + this.volume.String()
	}
	if this.app != nil && this.app.AppId != "" {
		str += " app=" + this.app.String()
	}

	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *device) setService(srv gopi.RPCServiceRecord) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.service = srv
	this.txt = nil
}

func (this *device) setStateVolume(v volume) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.volume == nil {
		this.volume = &v
	} else if this.volume.Equals(v) == false {
		this.volume = &v
		this.Log.Info("Volume changed=", v) // TODO EMIT
	}
}

func (this *device) setStateApplication(app application) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.app == nil {
		this.app = &app
	} else if this.app.Equals(app) == false {
		this.app = &app
		this.Log.Info("App changed=", app) // TODO EMIT
	}
}

func (this *device) Txt(key string) string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.txt == nil {
		this.txt = make(map[string]string)
		for _, txt := range this.service.Txt {
			if pair := strings.SplitN(txt, "=", 2); len(pair) == 2 {
				this.txt[pair[0]] = pair[1]
			}
		}
	}
	if value, exists := this.txt[key]; exists {
		return value
	} else {
		return ""
	}
}
