/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type connection struct {
	conn *tls.Conn
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT AND DISCONNECT

func (this *connection) Connect(srv gopi.RPCServiceRecord, flags gopi.RPCFlag, timeout time.Duration) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// If already connected, return
	if this.conn != nil {
		return gopi.ErrOutOfOrder
	}

	// Get a connection address and port
	if addr, port, err := getAddrPort(srv, flags); err != nil {
		return err
	} else if conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}, "tcp", fmt.Sprintf("%s:%d", addr, port), &tls.Config{
		InsecureSkipVerify: true,
	}); err != nil {
		return fmt.Errorf("%s: %w", addr, err)
	} else {
		this.conn = conn
	}

	// Success
	return nil
}

func (this *connection) Disconnect() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// If already connected, return
	if this.conn == nil {
		return nil
	}

	// Close connection
	err := this.conn.Close()
	this.conn = nil
	return err
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *connection) LocalAddr() string {
	if this.conn != nil {
		return this.conn.LocalAddr().String()
	} else {
		return "<nil>"
	}
}

func (this *connection) RemoteAddr() string {
	if this.conn != nil {
		return this.conn.RemoteAddr().String()
	} else {
		return "<nil>"
	}
}

func (this *connection) IsConnected() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.conn != nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *connection) String() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return "<cast.Connection addr=" + this.RemoteAddr() + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func getAddrPort(srv gopi.RPCServiceRecord, flag gopi.RPCFlag) (net.IP, uint, error) {
	// Set default flags
	if flag&(gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6) == 0 {
		flag |= (gopi.RPC_FLAG_INET_V4 | gopi.RPC_FLAG_INET_V6)
	}

	// Make pool of addresses
	addrs := make([]net.IP, 0, len(srv.Addrs))
	for _, addr := range srv.Addrs {
		if addr.To16() != nil && (flag&gopi.RPC_FLAG_INET_V6) != 0 {
			addrs = append(addrs, addr)
		} else if addr.To4() != nil && (flag&gopi.RPC_FLAG_INET_V4) != 0 {
			addrs = append(addrs, addr)
		}
	}

	// Check we have valid addrs and port
	if len(addrs) == 0 || srv.Port == 0 {
		return nil, 0, gopi.ErrBadParameter
	}

	// If RPC_FLAG_SERVICE_ANY return a random address
	// or else return the first one
	index := 0
	if flag&gopi.RPC_FLAG_SERVICE_ANY != 0 {
		index = rand.Intn(len(addrs) - 1)
	}

	// Return address and port
	return addrs[index], uint(srv.Port), nil
}
