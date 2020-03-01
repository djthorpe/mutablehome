/*
	Mutablehome Automation: Web Server
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package httpd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Httpd struct {
	Iface net.Interface
	Port  uint
	Flags gopi.RPCFlag
}

type httpd struct {
	log    gopi.Logger
	server *http.Server
	iface  net.Interface
	host   string
	port   uint

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SERVICE_TYPE_HTTP = "_http._tcp"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Httpd) Name() string { return "httpd" }

func (config Httpd) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(httpd)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *httpd) Init(config Httpd) error {

	// Set port
	if config.Port > 0 {
		this.port = config.Port
	} else if port, err := unusedPort(); err != nil {
		return err
	} else {
		this.port = port
	}

	// Ensure one or more addresses
	if config.Iface.Index > 0 {
		if config.Iface.Flags&net.FlagUp != net.FlagUp {
			return gopi.ErrBadParameter.WithPrefix("Interface down")
		}
		if ip, err := addrForInterface(config.Iface, config.Flags); err != nil {
			return err
		} else if ip.To4() != nil {
			this.host = ip.String()
		} else if ip.To16() != nil {
			this.host = fmt.Sprintf("[%v]", ip.String())
		}
		this.iface = config.Iface
	} else if host, err := os.Hostname(); err != nil {
		return err
	} else {
		this.host = host
	}

	// Return success
	return nil
}

func (this *httpd) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.server != nil {
		if err := this.server.Close(); err != nil {
			return err
		}
	}

	// Return success
	return this.Unit.Close()
}

func addrForInterface(iface net.Interface, flags gopi.RPCFlag) (net.IP, error) {
	if flags&(gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6) == 0 {
		flags = gopi.RPC_FLAG_INET_V4 | gopi.RPC_FLAG_INET_V6
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if addr, _, err := net.ParseCIDR(addr.String()); err == nil {
			if addr.To16() != nil && flags&gopi.RPC_FLAG_INET_V6 == gopi.RPC_FLAG_INET_V6 {
				return addr, nil
			} else if addr.To4() != nil && flags&gopi.RPC_FLAG_INET_V4 == gopi.RPC_FLAG_INET_V4 {
				return addr, nil
			}
		}
	}
	return nil, gopi.ErrNotFound
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *httpd) String() string {
	return "<" + this.Log.Name() +
		" addr=" + this.Addr() +
		" host=" + strconv.Quote(this.host) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION HttpServer

func (this *httpd) Addr() string {
	if this.iface.Index == 0 {
		return fmt.Sprintf(":%v", this.port)
	} else {
		return fmt.Sprintf("%v:%v", this.host, this.port)
	}
}

func (this *httpd) HostPort() string {
	return fmt.Sprintf("%v:%v", this.host, this.port)
}

func (this *httpd) BaseURL() string {
	return fmt.Sprintf("http://%v/", this.HostPort())
}

func (this *httpd) ServeStatic(folder string) (*url.URL, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check folder
	if stat, err := os.Stat(folder); err != nil {
		return nil, err
	} else if stat.IsDir() == false {
		return nil, gopi.ErrBadParameter.WithPrefix(folder)
	} else {
		folder = filepath.Clean(folder)
	}

	// Make URL
	if url, err := url.Parse(this.BaseURL()); err != nil {
		return nil, err
	} else if this.server != nil {
		return nil, gopi.ErrOutOfOrder
	} else {
		this.server = &http.Server{}
		this.server.Handler = http.FileServer(http.Dir(folder))
		this.server.Addr = this.HostPort()
		go func() {
			this.server.ListenAndServe()
		}()
		return url, nil
	}
}

func (this *httpd) Stop(ctx context.Context) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.server != nil {
		return this.server.Shutdown(ctx)
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func unusedPort() (uint, error) {
	if addr, err := net.ResolveTCPAddr("tcp", ":0"); err != nil {
		return 0, err
	} else if sock, err := net.ListenTCP("tcp", addr); err != nil {
		return 0, err
	} else {
		defer sock.Close()
		return uint(sock.Addr().(*net.TCPAddr).Port), nil
	}
}
