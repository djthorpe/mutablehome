/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"context"
	"fmt"
	"sync"
	"time"

	// Frameworks
	grpc "github.com/djthorpe/gopi-rpc/v2/unit/grpc"
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	"github.com/djthorpe/mutablehome"

	// Protocol buffers
	pb "github.com/djthorpe/mutablehome/protobuf/mutablehome"
	ptypes "github.com/golang/protobuf/ptypes"
	empty "github.com/golang/protobuf/ptypes/empty"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type NodeService struct {
	Server gopi.RPCServer
}

type nodeservice struct {
	base.Unit
	sync.Mutex

	server gopi.RPCServer
	node   mutablehome.Node
	start  time.Time
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (NodeService) Name() string { return "rpc/mutablehome/node" }

func (config NodeService) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(nodeservice)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else if err := this.Init(config); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *nodeservice) Init(config NodeService) error {
	// Set server
	if config.Server == nil {
		return gopi.ErrBadParameter.WithPrefix("Server")
	} else {
		this.server = config.Server
	}

	// Register with server
	pb.RegisterNodeServer(this.server.(grpc.GRPCServer).GRPCServer(), this)

	// Set start time
	this.start = time.Now()

	// Success
	return nil
}

func (this *nodeservice) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Release resources
	this.server = nil
	this.node = nil

	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *nodeservice) String() string {
	return "<" + this.Log.Name() + " " + fmt.Sprint(this.server) + " node=" + fmt.Sprint(this.node) + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.RPCNodeService

func (this *nodeservice) SetNode(node mutablehome.Node) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if node == nil || this.node != nil {
		return gopi.ErrBadParameter.WithPrefix("node")
	} else {
		this.node = node
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCService

func (this *nodeservice) CancelRequests() error {
	// Do not need to cancel
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *nodeservice) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	this.Log.Debug("<Ping>")

	return &empty.Empty{}, nil
}

func (this *nodeservice) Metadata(context.Context, *empty.Empty) (*pb.MetadataResponse, error) {
	this.Log.Debug("<Metadata>")
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check to make sure node is set
	if this.node == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Missing node parameter")
	}

	// Return metadata information
	return &pb.MetadataResponse{
		Id:     this.node.Id(),
		Name:   this.node.Name(),
		Uptime: ptypes.DurationProto(time.Now().Sub(this.start)),
	}, nil
}
