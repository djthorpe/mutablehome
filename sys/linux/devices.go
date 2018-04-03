/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Devices struct {
	// The root path for all mutablehome data
	Root string

	// The filename for the device database
	Filename string
}

type devices struct {
	log  gopi.Logger
	path string
	arr  []*mutablehome.Device
	hash map[string]*mutablehome.Device
	lock sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_ROOT             = ".mutablehome"
	DEFAULT_FILE_PERMISSIONS = 0644
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Devices) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<mutablehome.sys.linux.devices>{ root=%v filename=%v }", config.Root, config.Filename)

	this := new(devices)
	this.log = log
	this.arr = make([]*mutablehome.Device, 0)
	this.hash = make(map[string]*mutablehome.Device)

	// If root folder is relative, make absolute based on home folder
	if config.Root == "" {
		config.Root = DEFAULT_ROOT
	}

	// Check for root path
	if root, exists := ResolvePath(config.Root, UserDir()); exists == false {
		return nil, fmt.Errorf("Path does not exist: %v", root)
	} else if IsWritableFolder(root) == false {
		return nil, fmt.Errorf("Path is not writable: %v", root)
	} else {
		this.path = filepath.Join(root, config.Filename)
	}

	// Read from disk - or write if it's not available
	/*
		if err := this.readfile(); os.IsNotExist(err) {
			if err := this.writefile(); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	*/

	// Return success
	return this, nil
}

func (this *devices) Close() error {
	this.log.Debug("<mutablehome.sys.linux.devices>{ path=%v }", this.path)

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

// Device returns an existing device or creates a new device
func (this *devices) Device(device_id uint64, device_type mutablehome.DeviceType, product_id uint64) (*mutablehome.Device, error) {
	if hash := mutablehome.Hash(device_type, device_id); hash == "" {
		this.log.Debug2("Device: invalid hash: device_type=%v device_id=%v", device_type, device_id)
		return nil, gopi.ErrBadParameter
	} else if device, exists := this.hash[hash]; exists {
		return device, nil
	} else {
		// Create a new device
		device := &mutablehome.Device{
			DeviceId:       device_id,
			Type:           device_type,
			ProductId:      product_id,
			Name:           "Unknown product",
			Location:       "Unknown location",
			Paired:         mutablehome.PAIR_STATUS_DISCOVERED,
			TimeDiscovered: time.Now(),
			TimeUnpaired:   time.Time{},
			TimePaired:     time.Time{},
			TimeUpdated:    time.Time{},
		}
		// Append the device
		this.arr = append(this.arr, device)
		this.hash[hash] = device
		// Write the data
		if err := this.writefile(); err != nil {
			this.log.Error("writefile: %v", err)
		}
		// Return the device
		return device, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE FILE

func (this *devices) readfile() error {

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Read
	if data, err := ioutil.ReadFile(this.path); err != nil {
		return err
	} else if err := json.Unmarshal(data, &this.arr); err != nil {
		return err
	}

	// Re-construct the hash table
	this.hash = make(map[string]*mutablehome.Device, len(this.arr))
	for _, device := range this.arr {
		if hash := mutablehome.Hash(device.Type, device.DeviceId); hash != "" {
			this.hash[hash] = device
		} else {
			return fmt.Errorf("Invalid hash: device_type=%v device_id=%v", device.Type, device.DeviceId)
		}
	}

	return nil
}

func (this *devices) writefile() error {

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Write
	if data, err := json.Marshal(this.arr); err != nil {
		return err
	} else if err := ioutil.WriteFile(this.path, data, DEFAULT_FILE_PERMISSIONS); err != nil {
		return err
	}

	// Success
	return nil
}
