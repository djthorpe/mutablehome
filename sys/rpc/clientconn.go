/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/djthorpe/mutablehome"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	evt "github.com/djthorpe/gopi/util/event"
	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
	reflection_pb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GRPCClientConn interface {
	mutablehome.RPCClientConn

	// Return the GRPC Connection
	Conn() *grpc.ClientConn
}

type ClientConn struct {
	Name       string
	Addr       string
	SSL        bool
	SkipVerify bool
	Timeout    time.Duration
}

type clientconn struct {
	log        gopi.Logger
	name       string
	addr       string
	ssl        bool
	skipverify bool
	timeout    time.Duration
	conn       *grpc.ClientConn
	pubsub     *evt.PubSub
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ClientConn) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<rpc.clientconn>Open(name=%v,addr=%v,ssl=%v,skipverify=%v,timeout=%v)", config.Name, config.Addr, config.SSL, config.SkipVerify, config.Timeout)

	// Create a client object
	this := new(clientconn)
	this.name = config.Name
	this.addr = config.Addr
	this.ssl = config.SSL
	this.skipverify = config.SkipVerify
	this.timeout = config.Timeout
	this.log = log
	this.conn = nil
	this.pubsub = evt.NewPubSub(0)

	// Success
	return this, nil
}

func (this *clientconn) Close() error {
	this.log.Debug("<rpc.clientconn>Close{ name=%v addr=%v }", this.name, this.addr)

	// Disconnect first
	err := this.Disconnect()

	// Then free any resources
	this.pubsub.Close()
	this.pubsub = nil

	// Return any error conditions
	return err
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT AND DISCONNECT

func (this *clientconn) Connect() error {
	this.log.Debug2("<rpc.clientconn>Connect{ name=%v addr=%v }", this.name, this.addr)
	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	if this.conn != nil {
		this.log.Debug("<rpc.clientconn>Connect: Cannot call Connect() when connection already made")
		return gopi.ErrOutOfOrder
	}

	// Create connection options
	opts := make([]grpc.DialOption, 0, 1)

	// SSL options
	if this.ssl {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: this.skipverify})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Connection timeout options
	if this.timeout > 0 {
		opts = append(opts, grpc.WithTimeout(this.timeout))
	}

	// Dial connection
	if conn, err := grpc.Dial(this.addr, opts...); err != nil {
		return err
	} else {
		this.conn = conn
	}

	// Emit a connected event
	this.emit(&clientevent{
		source: this,
		type_:  gopi.RPC_EVENT_CLIENT_CONNECTED,
	})

	return nil
}

func (this *clientconn) Disconnect() error {
	this.log.Debug2("<rpc.clientconn>Disconnect{ name=%v addr=%v }", this.name, this.addr)
	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil

		// Emit a disconnected event
		this.emit(&clientevent{
			source: this,
			type_:  gopi.RPC_EVENT_CLIENT_DISCONNECTED,
		})

		return err
	}
	return nil
}

func (this *clientconn) Services() ([]string, error) {
	this.log.Debug2("<rpc.clientconn>Services{}")
	if this.conn == nil {
		return nil, gopi.ErrOutOfOrder
	}
	reflection, err := this.newServerReflectionClient()
	if err != nil {
		return nil, err
	}
	defer reflection.CloseSend()
	if services, err := this.listServices(reflection); err != nil {
		return nil, err
	} else {
		return services, nil
	}
}

func (this *clientconn) Name() string {
	return this.name
}

func (this *clientconn) Addr() string {
	return this.addr
}

func (this *clientconn) Connected() bool {
	return this.conn != nil
}

func (this *clientconn) Conn() *grpc.ClientConn {
	return this.conn
}

////////////////////////////////////////////////////////////////////////////////
// SUBSCRIBE AND UNSUBSCRIBE

func (this *clientconn) Subscribe() <-chan gopi.Event {
	return this.pubsub.Subscribe()
}

func (this *clientconn) Unsubscribe(evt <-chan gopi.Event) {
	this.pubsub.Unsubscribe(evt)
}

func (this *clientconn) emit(evt gopi.RPCEvent) {
	this.pubsub.Emit(evt)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *clientconn) String() string {
	return fmt.Sprintf("<rpc.ClientConn>{ name=%v addr=%v ssl=%v connected=%v }", this.name, this.addr, this.ssl, this.conn != nil)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *clientconn) newServerReflectionClient() (reflection_pb.ServerReflection_ServerReflectionInfoClient, error) {
	if this.conn == nil {
		return nil, gopi.ErrOutOfOrder
	}
	ctx := context.Background()
	if this.timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, this.timeout)
	}
	if client, err := reflection_pb.NewServerReflectionClient(this.conn).ServerReflectionInfo(ctx); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (this *clientconn) listServices(c reflection_pb.ServerReflection_ServerReflectionInfoClient) ([]string, error) {
	if err := c.Send(&reflection_pb.ServerReflectionRequest{
		MessageRequest: &reflection_pb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return nil, err
	}
	if resp, err := c.Recv(); err != nil {
		return nil, err
	} else if modules := resp.GetListServicesResponse(); modules == nil {
		return nil, fmt.Errorf("GetListServicesResponse() error")
	} else {
		module_services := modules.GetService()
		module_names := make([]string, len(module_services))
		for i, service := range module_services {
			// Full name of a registered service, including its package name
			module_names[i] = service.Name
		}
		return module_names, nil
	}
}
