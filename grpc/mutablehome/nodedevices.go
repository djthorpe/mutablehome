/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type nodedevices struct {
	sync.Mutex
	sync.WaitGroup

	Log     gopi.Logger
	online  bool
	node    mutablehome.Node
	stop    chan struct{}
	devices map[string]bool
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *nodedevices) Init(log gopi.Logger) error {
	// Set logger
	if log == nil {
		return gopi.ErrBadParameter.WithPrefix("Log")
	} else {
		this.Log = log
	}

	// Create stop channel
	this.stop = make(chan struct{})

	// Success
	return nil
}

func (this *nodedevices) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Send stop signal and wait for stop
	close(this.stop)
	this.WaitGroup.Wait()

	// Release resources
	this.stop = nil
	this.node = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.RPCNodeService

func (this *nodedevices) SetNode(node mutablehome.Node) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if node == nil || this.node != nil {
		return gopi.ErrBadParameter.WithPrefix("node")
	} else {
		this.node = node
		go this.BackgroundProcess(this.node, this.stop)
	}

	// Return success
	return nil
}

func (this *nodedevices) SetOnline(value bool) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if this.node != nil {
		this.online = value
		if value {
			this.Log.Info(this.node.Id(), ": Online")
		} else {
			this.Log.Info(this.node.Id(), ": Offline")
		}
	} else {
		this.online = false
	}
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND PROCESS

func (this *nodedevices) BackgroundProcess(node mutablehome.Node, stop <-chan struct{}) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

	evts := node.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-evts:
			if evt_, ok := evt.(mutablehome.Event); ok {
				this.ProcessEvent(evt_)
			}
		case <-stop:
			node.Unsubscribe(evts)
			break FOR_LOOP
		}
	}
}

func (this *nodedevices) ProcessEvent(evt mutablehome.Event) {
	switch evt.Type() {
	case mutablehome.EVENT_NODE_ONLINE:
		this.SetOnline(true)
	case mutablehome.EVENT_NODE_OFFLINE:
		this.SetOnline(false)
	case mutablehome.EVENT_DEVICE_ADDED:
		if err := this.AddDevice(evt.Device()); err != nil {
			this.Log.Error(fmt.Errorf("Ignoring: %v: %w", evt.Type(), err))
		}
	case mutablehome.EVENT_DEVICE_REMOVED:
		this.Log.Warn("Ignoring:", evt)
	case mutablehome.EVENT_DEVICE_CHANGED:
		this.Log.Warn("Ignoring:", evt)
	default:
		this.Log.Warn("Ignoring:", evt)
	}
}

////////////////////////////////////////////////////////////////////////////////
// DEVICE METHODS

func (this *nodedevices) AddDevice(device mutablehome.Device) error {
	if device == nil {
		return gopi.ErrBadParameter.WithPrefix("Device")
	}

	this.Log.Debug("TODO:", device)

	// Successs
	return nil
}
