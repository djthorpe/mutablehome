/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	// Frameworks
	app "github.com/djthorpe/gopi-rpc/v2/app"
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi-rpc/v2/unit/grpc"
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
	_ "github.com/djthorpe/mutablehome/grpc/mutablehome"
	_ "github.com/djthorpe/mutablehome/unit/tradfri"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	TRADFRI_MDNS_SERVICE = "_coap._udp"
	TRADFRI_MDNS_TIMEOUT = 2 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// SERVICE DISCOVERY

func Services(app gopi.App) ([]gopi.RPCServiceRecord, error) {
	// Lookup the gateway
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	ctx, cancel := context.WithTimeout(context.Background(), TRADFRI_MDNS_TIMEOUT)
	defer cancel()
	return discovery.Lookup(ctx, TRADFRI_MDNS_SERVICE)
}

////////////////////////////////////////////////////////////////////////////////
// OBSERVE DEVICES

func ObserveDevices(app gopi.App, tradfri mutablehome.Ikea) (*sync.WaitGroup, []context.CancelFunc, error) {
	wg := new(sync.WaitGroup)
	if devices, err := tradfri.Devices(); err != nil {
		return nil, nil, err
	} else {
		cancels := make([]context.CancelFunc, 0, len(devices))
		for _, device := range devices {
			ctx, cancel := context.WithCancel(context.Background())
			wg.Add(1)
			go func(device uint) {
				if err := tradfri.ObserveDevice(ctx, device); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
					app.Log().Error(err)
				}
				wg.Done()
			}(device)
			cancels = append(cancels, cancel)
		}
		return wg, cancels, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// CONNECT METHODS

func ConnectService(app gopi.App, tradfri mutablehome.Ikea, service gopi.RPCServiceRecord) error {
	if err := tradfri.Connect(service, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		return err
	} else {
		return nil
	}
}

func ConnectHostPort(_ gopi.App, tradfri mutablehome.Ikea, host, port string) error {
	return gopi.ErrNotImplemented
}

func ConnectNameAddr(app gopi.App, tradfri mutablehome.Ikea, addr string) error {
	matched := make([]gopi.RPCServiceRecord, 0, 1)
	if services, err := Services(app); err != nil {
		return err
	} else if len(services) == 0 {
		return fmt.Errorf("No Tradfri gateway found")
	} else {
		for _, service := range services {
			if service.Name == addr {
				matched = append(matched, service)
			} else if service.Host == addr {
				matched = append(matched, service)
			}
		}
	}

	if len(matched) == 0 {
		return fmt.Errorf("No Tradfri gateway found")
	} else if len(matched) > 1 {
		return fmt.Errorf("More than one Tradfri gateway found, use -addr to select between them")
	} else {
		return ConnectService(app, tradfri, matched[0])
	}
}

func Connect(app gopi.App, tradfri mutablehome.Ikea) error {
	if services, err := Services(app); err != nil {
		return err
	} else if len(services) == 0 {
		return fmt.Errorf("No Tradfri gateway found")
	} else if len(services) > 1 {
		return fmt.Errorf("More than one Tradfri gateway found, use -addr to select between them")
	} else {
		return ConnectService(app, tradfri, services[0])
	}
}

////////////////////////////////////////////////////////////////////////////////
// MAIN

func Main(app gopi.App, args []string) error {
	// Don't allow any arguments
	if len(args) != 0 {
		return fmt.Errorf("Arguments provided but not required")
	}

	// Tradfri unit
	tradfri := app.UnitInstance("tradfri").(mutablehome.Ikea)

	// Connect to Tradfri
	if addr := app.Flags().GetString("addr", gopi.FLAG_NS_DEFAULT); addr != "" {
		if host, port, err := net.SplitHostPort(addr); err == nil {
			if err := ConnectHostPort(app, tradfri, host, port); err != nil {
				return err
			}
		} else if err := ConnectNameAddr(app, tradfri, addr); err != nil {
			return err
		}
	} else if err := Connect(app, tradfri); err != nil {
		return err
	}

	// Observe all the devices so we can receive events from them
	if wg, cancels, err := ObserveDevices(app, tradfri); err != nil {
		return err
	} else {
		// Wait until CTRL+C pressed
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal(context.Background(), os.Interrupt)

		// Cancel observing devices and wait until done
		for _, cancel := range cancels {
			cancel()
		}
		wg.Wait()
	}

	// Disconnect
	tradfri.Disconnect()

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewServer(Main, "rpc/mutablehome/node", "mutablehome/node/tradfri", "register", "discovery"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// -addr is the address to a tradfri gateway
		app.Flags().FlagString("addr", "", "Tradfri gateway address")
		// Run and exit
		os.Exit(app.Run())
	}
}
