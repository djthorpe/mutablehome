/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	iface "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cast struct {
	Discovery gopi.RPCServiceDiscovery
}

type cast struct {
	log       gopi.Logger
	discovery gopi.RPCServiceDiscovery

	Lookup
	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// COMSTANTS

const (
	SERVICE_TYPE_GOOGLECAST = "_googlecast._tcp"
	DELTA_LOOKUP_TIME       = 60 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Cast) Name() string { return "googlecast" }

func (config Cast) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(cast)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *cast) Init(config Cast) error {
	// Check for discovery
	if config.Discovery == nil {
		return gopi.ErrBadParameter.WithPrefix("discovery")
	} else {
		this.Lookup.Discovery = config.Discovery
	}

	// Start discovery
	if err := this.Lookup.Start(SERVICE_TYPE_GOOGLECAST); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *cast) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Stop discovery
	this.Lookup.Stop()

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION Cast

func (this *cast) Devices() []iface.Device {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return nil
}
