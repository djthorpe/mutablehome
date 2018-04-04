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
	"math/rand"
	"os"
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
	log          gopi.Logger
	path         string
	arr          []*mutablehome.Device
	hash         map[string]*mutablehome.Device
	lock         sync.Mutex
	ticker_delta time.Duration
	ticker_quit  chan struct{}
	ticker_done  chan struct{}
	ticker_sync  bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DEFAULT_ROOT             = ".mutablehome"
	DEFAULT_DELTA_SECS       = 5
	DEFAULT_FILE_PERMISSIONS = 0644
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Devices) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<mutablehome.sys.linux.devices>Open{ root=%v filename=%v }", config.Root, config.Filename)

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
	if err := this.readfile(); os.IsNotExist(err) {
		if err := this.writefile(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Set up the ticker which syncs writing of the devices
	// to a file occasionally
	this.ticker_delta = time.Second * DEFAULT_DELTA_SECS
	this.ticker_quit = make(chan struct{})
	this.ticker_done = make(chan struct{})
	this.ticker_sync = false
	go this.ticker()

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Return success
	return this, nil
}

func (this *devices) Close() error {
	this.log.Debug("<mutablehome.sys.linux.devices>Close{ path=%v }", this.path)

	// Sync write and quit ticker
	this.ticker_quit <- gopi.DONE
	<-this.ticker_done

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

// Device returns an existing device or creates a new device
func (this *devices) Device(device_id uint64, device_type mutablehome.DeviceType, product_id uint64) (*mutablehome.Device, error) {
	this.log.Debug2("<mutablehome.sys.linux.devices>Device{ device_id=%X device_type=%v product_id=%X }", device_id, device_type, product_id)

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
			PairStatus:     mutablehome.PAIR_STATUS_DISCOVERED,
			TimeDiscovered: time.Now(),
			TimeUnpaired:   time.Time{},
			TimePaired:     time.Time{},
			TimeUpdated:    time.Time{},
		}

		// Append the device
		this.arr = append(this.arr, device)
		this.hash[hash] = device

		// Sync
		this.sync(true)

		// Return the device
		return device, nil
	}
}

// Devices returns a set of devices in the database
func (this *devices) Devices(device_type_flag mutablehome.DeviceType, pair_status_flag mutablehome.PairStatusType) []*mutablehome.Device {
	this.log.Debug2("<mutablehome.sys.linux.devices>Devices{ device_type_flag=%v pair_status_flag=%v }", device_type_flag, pair_status_flag)

	devices := make([]*mutablehome.Device, 0)
	for _, device := range this.arr {
		if device.Type&device_type_flag == 0 {
			continue
		} else if device.PairStatus&pair_status_flag == 0 {
			continue
		} else {
			devices = append(devices, device)
		}
	}

	// Append a new (fake) Energenie Control device which can be used for pairing
	if (device_type_flag&mutablehome.DEVICE_TYPE_ENERGENIE_CONTROL) != 0 && (pair_status_flag&mutablehome.PAIR_STATUS_DISCOVERED) != 0 {
		if new_device := this.newEnergenieControlDevice(); new_device != nil {
			devices = append(devices, new_device)
		} else {
			this.log.Warn("Unable to append Energenie Control device, ignoring")
		}
	}

	return devices
}

// Pair a device
func (this *devices) Pair(device_id uint64, device_type mutablehome.DeviceType) error {
	this.log.Debug2("<mutablehome.sys.linux.devices>Pair{ device_id=%X device_type=%v }", device_id, device_type)

	if hash := mutablehome.Hash(device_type, device_id); hash == "" {
		this.log.Debug2("Pair: invalid hash: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrBadParameter
	} else if device, exists := this.hash[hash]; exists == false {
		this.log.Debug2("Pair: not found: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrBadParameter
	} else if device.PairStatus != mutablehome.PAIR_STATUS_DISCOVERED || device.PairStatus != mutablehome.PAIR_STATUS_UNPAIRED {
		this.log.Debug2("Pair: invalid state: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrOutOfOrder
	} else {
		// Set device as paired
		device.PairStatus = mutablehome.PAIR_STATUS_PAIRED
		device.TimeUnpaired = time.Time{}
		device.TimePaired = time.Now()

		// Sync
		this.sync(true)

		// TODO: Emit a "pair" event
	}

	return nil
}

// Unpair a device
func (this *devices) Unpair(device_id uint64, device_type mutablehome.DeviceType) error {
	this.log.Debug2("<mutablehome.sys.linux.devices>Unpair{ device_id=%X device_type=%v }", device_id, device_type)

	if hash := mutablehome.Hash(device_type, device_id); hash == "" {
		this.log.Debug2("Unpair: invalid hash: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrBadParameter
	} else if device, exists := this.hash[hash]; exists == false {
		this.log.Debug2("Unpair: not found: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrBadParameter
	} else if device.PairStatus != mutablehome.PAIR_STATUS_DISCOVERED || device.PairStatus != mutablehome.PAIR_STATUS_PAIRED {
		this.log.Debug2("Unpair: invalid state: device_type=%v device_id=%v", device_type, device_id)
		return gopi.ErrOutOfOrder
	} else {
		// Set device as unpaired
		device.PairStatus = mutablehome.PAIR_STATUS_UNPAIRED
		device.TimeUnpaired = time.Now()

		// Sync
		this.sync(true)

		// TODO: Emit an "unpair" event
	}

	return nil
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

func (this *devices) sync(flag bool) {
	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()
	// Flag that we should sync or clear flag
	this.ticker_sync = flag
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

func (this *devices) ticker() {
	ticker := time.NewTicker(this.ticker_delta)

FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			if this.ticker_sync {
				if err := this.writefile(); err != nil {
					this.log.Error("writefile: %v", err)
				}
				this.sync(false)
			}
		case <-this.ticker_quit:
			break FOR_LOOP
		}
	}

	// Sync write
	if this.ticker_sync {
		if err := this.writefile(); err != nil {
			this.log.Error("writefile: %v", err)
		}
		this.sync(false)
	}

	// Indicate done
	this.ticker_done <- gopi.DONE
}

////////////////////////////////////////////////////////////////////////////////
// GENERATE ENERGENIE DEVICE

// newEnergenieControlId attempts to create a non-existing Sensor ID over 100
// iterations or else returns 0
func (this *devices) newEnergenieControlId() uint64 {
	for i := 0; i < 100; i++ {
		// 20-bit random number which is non-zero
		device_id := rand.Uint64() & 0xFFFFF
		// Check for nil or collisions
		if device_id == 0 {
			continue
		} else if hash := mutablehome.Hash(mutablehome.DEVICE_TYPE_ENERGENIE_CONTROL, device_id); hash == "" {
			continue
		} else if _, exists := this.hash[hash]; exists {
			continue
		} else {
			return device_id
		}
	}
	// Failed to find a valid id
	return 0
}

// newEnergenieControlDevice attempts to create a non-existing
// device structure or else returns nil
func (this *devices) newEnergenieControlDevice() *mutablehome.Device {
	if device_id := this.newEnergenieControlId(); device_id == 0 {
		return nil
	} else {
		return &mutablehome.Device{
			DeviceId:   device_id,
			Type:       mutablehome.DEVICE_TYPE_ENERGENIE_CONTROL,
			Name:       "New Energenie Control Device",
			PairStatus: mutablehome.PAIR_STATUS_DISCOVERED,
		}
	}
}
