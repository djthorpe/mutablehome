/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type devices struct {
	devices map[string]*device

	sync.Mutex
}

type device struct {
	gopi.RPCServiceRecord
	txt_ map[string]string

	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// ADD, REMOVE AND UPDATE DEVICES

func (this *devices) Init() {
	this.devices = make(map[string]*device)
}

func (this *devices) Close() {
	this.devices = nil
}

func (this *devices) Update(srv gopi.RPCServiceRecord) (*device, bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if key := srv.Name; key == "" {
		return nil, false
	} else if _, exists := this.devices[key]; exists {
		this.devices[key].Lock()
		defer this.devices[key].Unlock()
		this.devices[key].txt_ = nil
		this.devices[key].RPCServiceRecord = srv
		return this.devices[key], true // TODO: Check equals
	} else {
		this.devices[key] = &device{RPCServiceRecord: srv}
		return this.devices[key], true
	}
}

func (this *devices) Remove(srv gopi.RPCServiceRecord) (*device, bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if key := srv.Name; key == "" {
		return nil, false
	} else if device, exists := this.devices[key]; exists == false {
		return nil, false
	} else {
		delete(this.devices, key)
		return device, true
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device

func (this *device) Id() string {
	return this.txt("id")
}

func (this *device) Name() string {
	return this.txt("fn")
}

func (this *device) Model() string {
	return this.txt("md")
}

func (this *device) Service() string {
	return this.txt("rs")
}

func (this *device) State() uint {
	if value := this.txt("st"); value == "" {
		return 0
	} else if value_, err := strconv.ParseUint(value, 10, 32); err != nil {
		return 0
	} else {
		return uint(value_)
	}
}

func (this *device) Equals(other *device) bool {
	if this.Id() != other.Id() {
		return false
	}
	if this.Name() != other.Name() {
		return false
	}
	if this.Model() != other.Model() {
		return false
	}
	if this.Service() != other.Service() {
		return false
	}
	if this.State() != other.State() {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	return "<cast.Device" +
		" id=" + strconv.Quote(this.Id()) +
		" name=" + strconv.Quote(this.Id()) +
		" model=" + strconv.Quote(this.Id()) +
		" service=" + strconv.Quote(this.Id()) +
		" state=" + fmt.Sprint(this.State()) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *device) txt(key string) string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.txt_ == nil {
		this.txt_ = make(map[string]string)
		for _, txt := range this.RPCServiceRecord.Txt {
			if pair := strings.SplitN(txt, "=", 2); len(pair) == 2 {
				this.txt_[pair[0]] = pair[1]
			}
		}
	}
	if value, exists := this.txt_[key]; exists {
		return value
	} else {
		return ""
	}
}
