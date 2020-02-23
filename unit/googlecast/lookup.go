/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"context"
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Lookup struct {
	Discovery gopi.RPCServiceDiscovery
	cancel    context.CancelFunc

	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// Lookup

// Start lookup in background
func (this *Lookup) Start(service string) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.cancel != nil {
		return gopi.ErrOutOfOrder
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		this.WaitGroup.Add(1)
		if services, err := this.Discovery.Lookup(ctx, service); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			fmt.Println(err)
		} else {
			fmt.Println(services)
		}
		this.WaitGroup.Done()
	}()
	this.cancel = cancel

	// Return success
	return nil
}

// Stop lookup
func (this *Lookup) Stop() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Cancel and wait until ended
	if this.cancel != nil {
		this.cancel()
		this.WaitGroup.Wait()
	}
}
