/*
	Mutablehome: Home Automation in Golang
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    https://github.com/djthorpe/mutablehome/
	For Licensing and Usage information, please see LICENSE
*/

package mutablehome

import (
	"context"
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"
	"github.com/djthorpe/gopi/util/event"

	// Protocol buffers
	pb "github.com/djthorpe/mutablehome/rpc/protobuf/mutablehome"
	empty "github.com/golang/protobuf/ptypes/empty"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MutableHome struct {
	Server    gopi.RPCServer
	Discovery gopi.RPCServiceDiscovery
}

type service struct {
	log       gopi.Logger
	discovery gopi.RPCServiceDiscovery

	event.Publisher
	event.Tasks
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config MutableHome) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.mutablehome>Open{ server=%v discovery=%v }", config.Server, config.Discovery)

	// Check for bad input parameters
	if config.Server == nil || config.Discovery == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.discovery = config.Discovery

	// Register service with GRPC server
	pb.RegisterMutableHomeServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Background task to connect & disconnect from mihome stubs
	this.Tasks.Start(this.ConnectTask, this.LookupTask)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mutablehome>Close>{}")

	// Close publisher
	this.Publisher.Close()

	// Stop background tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("<grpc.service.mutablehome>{ }")
}

////////////////////////////////////////////////////////////////////////////////
// CANCEL STREAMING REQUESTS

func (this *service) CancelRequests() error {
	this.log.Debug2("<grpc.service.mutablehome>CancelRequests{}")

	// Cancel any streaming requests
	this.Publisher.Emit(event.NullEvent)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RPC METHODS

// Ping returns an empty response
func (this *service) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.log.Debug("<grpc.service.mutablehome>Ping{ }")

	this.Lock()
	defer this.Unlock()

	return &empty.Empty{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *service) ConnectTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	events := this.discovery.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-events:
			if evt == nil {
			} else if evt_, ok := evt.(gopi.RPCEvent); ok {
				r := evt_.ServiceRecord()
				if r != nil && r.Service() == "_gopi._tcp" && r.Subtype() != "" {
					this.log.Debug("Event: %v: %v: %v:%v", evt_.Type(), r.Subtype(), r.Host(), r.Port())
				}
			}
		case <-stop:
			break FOR_LOOP
		}
	}

	this.discovery.Unsubscribe(events)

	// Success
	return nil
}

func (this *service) LookupTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	timer := time.NewTimer(1 * time.Second)
FOR_LOOP:
	for {
		select {
		case <-timer.C:
			// Lookup records
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			// Perform lookup
			if _, err := this.discovery.Lookup(ctx, "_gopi._tcp"); err != nil {
				this.log.Warn("Lookup: %v", err)
			}
			// Check again after 1 minute
			timer.Reset(1 * time.Minute)
		case <-stop:
			break FOR_LOOP
		}
	}

	// Success
	return nil
}
