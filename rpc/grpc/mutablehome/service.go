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
	"sync"

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
	Server gopi.RPCServer
}

type service struct {
	log gopi.Logger

	// Lock
	sync.Mutex

	// Emit events
	event.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config MutableHome) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.mutablehome>Open{ server=%v }", config.Server)

	// Check for bad input parameters
	if config.Server == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log

	// Register service with GRPC server
	pb.RegisterMutableHomeServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Background task to connect & disconnect from stubs

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.mutablehome>Close>{}")

	// Close publisher
	this.Publisher.Close()

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
