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
}

////////////////////////////////////////////////////////////////////////////////
// ADD, REMOVE AND UPDATE DEVICES

func (this *devices) Init() {
	this.devices = make(map[string]*device)
}

func (this *devices) Close() {
	this.devices = nil
}

func (this *devices) Update(srv gopi.RPCServiceRecord) bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if key := srv.Name; key == "" {
		return false
	} else {
		this.devices[key] = &device{srv}
		return true
	}
}

func (this *devices) Remove(srv gopi.RPCServiceRecord) bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if key := srv.Name; key == "" {
		return false
	} else if _, exists := this.devices[key]; exists == false {
		return false
	} else {
		delete(this.devices, key)
		return true
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION cast.Device

func (this *device) Id() string {
	return this.RPCServiceRecord.Name
}

func (this *device) Name() string {
	return this.RPCServiceRecord.Name
}

func (this *device) Model() string {
	return "MODEL"
}

func (this *device) Service() string {
	return "SERVICE"
}

func (this *device) State() uint {
	return 0
}

func (this *device) String() string {
	return "<cast.Device" +
		" id=" + strconv.Quote(this.Id()) +
		" name=" + strconv.Quote(this.Id()) +
		" model=" + strconv.Quote(this.Id()) +
		" service=" + strconv.Quote(this.Id()) +
		" state=" + fmt.Sprint(this.State()) +
		">"
}
