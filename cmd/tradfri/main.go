package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

const (
	TRADFRI_MDNS_SERVICE = "_coap._udp"
)

var (
	reHostPort = regexp.MustCompile("^([^\\:]+)\\:(\\d+)$")
	wg         sync.WaitGroup
)

/////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	tradfri := app.UnitInstance("tradfri").(mutablehome.Ikea)

	if services, err := Services(app); err != nil {
		return err
	} else if len(services) == 0 {
		return gopi.ErrNotFound.WithPrefix("Gateway")
	} else if len(services) != 1 {
		return gopi.ErrDuplicateItem.WithPrefix("Gateway")
	} else if err := tradfri.Connect(services[0], gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		return err
	}

	// Start observing devices
	cancels, err := ObserveDevices(tradfri)
	if err != nil {
		return err
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Cancel observing devices
	for _, cancel := range cancels {
		cancel()
	}

	// Wait for all goroutines to have ended
	wg.Wait()

	// Return success
	return nil
}

func Services(app gopi.App) ([]gopi.RPCServiceRecord, error) {
	addr := app.Flags().GetString("tradfri.addr", gopi.FLAG_NS_DEFAULT)
	// Deal with the case where tradfri.addr is given
	if addr != "" {
		if host, port, err := HostPort(app); err != nil {
			return nil, err
		} else if addrs, err := LookupHost(host); err != nil {
			return nil, err
		} else {
			return []gopi.RPCServiceRecord{
				{Addrs: addrs, Port: uint16(port)},
			}, nil
		}
	}
	// Lookup the gateway
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return discovery.Lookup(ctx, TRADFRI_MDNS_SERVICE)
}

func HostPort(app gopi.App) (string, uint, error) {
	addr := app.Flags().GetString("tradfri.addr", gopi.FLAG_NS_DEFAULT)

	// Add the port on the end if not added
	if addr == "" {
		return "", 0, gopi.ErrBadParameter.WithPrefix("-tradfri.addr")
	} else if reHostPort.MatchString(addr) == false {
		addr = fmt.Sprintf("%v:%v", addr, mutablehome.IKEA_DEFAULT_PORT)
	}

	// Check host and port
	if host, port, err := net.SplitHostPort(addr); err != nil {
		return "", 0, gopi.ErrBadParameter.WithPrefix("-tradfri.addr")
	} else if port_, err := strconv.ParseUint(port, 10, 32); err != nil {
		return "", 0, gopi.ErrBadParameter.WithPrefix("-tradfri.addr")
	} else {
		return host, uint(port_), nil
	}
}

func LookupHost(host string) ([]net.IP, error) {
	if addrs, err := net.LookupHost(host); err != nil {
		return nil, err
	} else {
		addrs_ := make([]net.IP, 0, len(addrs))
		for _, addr := range addrs {
			addrs_ = append(addrs_, net.ParseIP(addr))
		}
		return addrs_, nil
	}
}

func ObserveDevices(tradfri mutablehome.Ikea) ([]context.CancelFunc, error) {
	if devices, err := tradfri.Devices(); err != nil {
		return nil, err
	} else {
		cancels := make([]context.CancelFunc, 0, len(devices))
		for _, device := range devices {
			ctx, cancel := context.WithCancel(context.Background())
			wg.Add(1)
			go func(device uint) {
				if err := tradfri.ObserveDevice(ctx, device); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
					fmt.Println("ObserveDevice:", err)
				}
				wg.Done()
			}(device)
			cancels = append(cancels, cancel)
		}
		return cancels, nil
	}
}
