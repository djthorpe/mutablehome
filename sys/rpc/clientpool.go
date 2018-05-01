/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	evt "github.com/djthorpe/gopi/util/event"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientPool struct {
	SkipVerify bool
	Timeout    time.Duration
}

type clientpool struct {
	log        gopi.Logger
	skipverify bool
	timeout    time.Duration
	merger     evt.EventMerger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config ClientPool) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<rpc.clientpool>Open{ timeout=%v skipverify=%v }", config.Timeout, config.SkipVerify)
	this := new(clientpool)
	this.log = log
	this.skipverify = config.SkipVerify
	this.timeout = config.Timeout
	this.merger = evt.NewEventMerger()

	// Success
	return this, nil
}

func (this *clientpool) Close() error {
	this.log.Debug("<rpc.clientpool>Close{}")

	// Release resources
	this.merger.Close()
	this.merger = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func (this *clientpool) Connect(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) (mutablehome.RPCClientConn, error) {
	this.log.Debug2("<rpc.clientpool>Connect{ service=%v flags=%v }", service, flags)

	// TODO
	ssl := false

	// Determine the address
	if addr := addressFor(service, flags); addr == "" {
		return nil, gopi.ErrBadParameter
	} else if clientconn_, err := gopi.Open(ClientConn{
		Name:       service.Name,
		Addr:       addr + ":" + fmt.Sprint(service.Port),
		SSL:        ssl,
		SkipVerify: this.skipverify,
		Timeout:    this.timeout,
	}, this.log); err != nil {
		return nil, err
	} else if clientconn, ok := clientconn_.(mutablehome.RPCClientConn); ok == false {
		return nil, gopi.ErrOutOfOrder
	} else {
		// subscribe to events from connection
		this.merger.Add(clientconn.Subscribe())

		// Do connection
		if err := clientconn.Connect(); err != nil {
			return nil, err
		}

		// return success
		return clientconn, nil
	}
}

func (this *clientpool) Disconnect(conn mutablehome.RPCClientConn) error {
	this.log.Debug2("<rpc.clientpool>Disconnect{ conn=%v }", conn)
	return conn.Disconnect()
}

////////////////////////////////////////////////////////////////////////////////
// PUBSUB

func (this *clientpool) Subscribe() <-chan gopi.Event {
	return this.merger.Subscribe()
}

func (this *clientpool) Unsubscribe(evt <-chan gopi.Event) {
	this.merger.Unsubscribe(evt)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func addressFor(service *gopi.RPCServiceRecord, flags gopi.RPCFlag) string {
	if flags&gopi.RPC_FLAG_INET_UDP != 0 {
		// We don't support UDP connections
		return ""
	} else if flags&gopi.RPC_FLAG_INET_V6 != 0 {
		if len(service.IP6) == 0 {
			return ""
		} else {
			return service.IP6[0].String()
		}
	} else if flags&gopi.RPC_FLAG_INET_V4 != 0 {
		if len(service.IP4) == 0 {
			return ""
		} else {
			return service.IP4[0].String()
		}
	} else {
		return service.Host
	}
}
