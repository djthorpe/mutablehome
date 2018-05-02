/*
   Go Language Raspberry Pi Interface
   (c) Copyright David Thorpe 2016-2018
   All Rights Reserved
   Documentation http://djthorpe.github.io/gopi/
   For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/mutablehome"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"
	_ "github.com/djthorpe/mutablehome/sys/rpc"

	// RPC Clients
	_ "github.com/djthorpe/mutablehome/sys/clients"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Discovery struct {
	service     string
	cancel_func context.CancelFunc
	connection  []*gopi.RPCClientConn
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	discovery = make(map[string]*Discovery)
	services  = make(chan *gopi.RPCServiceRecord)
)

////////////////////////////////////////////////////////////////////////////////
// SERVICE DISCOVERY

func StopDiscovery(app *gopi.AppInstance, service_type string) error {
	if d, exists := discovery[service_type]; exists {
		app.Logger.Debug("StopDiscovery for '%v'", service_type)
		d.cancel_func()
	}
	return nil
}

func StartDiscovery(app *gopi.AppInstance, service string) error {
	mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	if service_type, err := gopi.RPCServiceType(service, gopi.RPC_FLAG_NONE); err != nil {
		return err
	} else if err := StopDiscovery(app, service_type); err != nil {
		return err
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		discovery[service_type] = &Discovery{cancel_func: cancel}
		go func() {
			app.Logger.Debug("Started discovery for service '%v'", service_type)
			if err := mdns.Browse(ctx, service_type); err != nil {
				app.Logger.Error("CancelDiscovery: %v", err)
			}
			app.Logger.Debug("Ended discovery for service '%v'", service_type)
			delete(discovery, service_type)
		}()
		return nil
	}
}

func DiscoveryLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	rpc_events := mdns.Subscribe()

FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-rpc_events:
			if evt == nil {
				continue
			} else if rpc_evt, ok := evt.(gopi.RPCEvent); ok == false {
				continue
			} else if rpc_evt.Type() != gopi.RPC_EVENT_SERVICE_RECORD {
				continue
			} else {
				services <- rpc_evt.ServiceRecord()
			}
		}
	}

	for service := range discovery {
		if err := StopDiscovery(app, service); err != nil {
			app.Logger.Error("CancelDiscovery: %v", err)
		}
	}

	mdns.Unsubscribe(rpc_events)
	return nil
}

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	// Subscribe to events from the client pool
	clientpool := app.ModuleInstance("rpc/clientpool").(mutablehome.RPCClientPool)
	poolevents := clientpool.Subscribe()

FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case service := <-services:
			if service == nil {
				continue
			} else if _, exists := discovery[service.Type]; exists == false {
				continue
			} else if err := Connect(app, service); err != nil {
				app.Logger.Error("Connect: %v: %v", service.Name, err)
			}
		case evt := <-poolevents:
			if err := HandleEvent(app, clientpool, evt); err != nil {
				app.Logger.Error("EventLoop: %v (for event %v)", err, evt)
			}
		}
	}

	// Unsubscribe from events from the client pool
	clientpool.Unsubscribe(poolevents)

	// Return success
	return nil
}

func HandleEvent(app *gopi.AppInstance, clientpool mutablehome.RPCClientPool, evt gopi.Event) error {
	if evt == nil {
		return nil
	} else if rpc_evt, ok := evt.(gopi.RPCEvent); ok == false {
		return nil
	} else if rpc_evt.Type() == gopi.RPC_EVENT_CLIENT_CONNECTED {
		// Obtain services for this connection
		conn := rpc_evt.Source().(mutablehome.RPCClientConn)
		if services, err := conn.Services(); err != nil {
			return err
		} else if HasService(services, "mutablelogic.Remotes") == false {
			return gopi.ErrNotImplemented
		} else if client, err := clientpool.NewClient("mutablelogic.Remotes", conn); err != nil {
			return err
		} else {
			codecs, err := client.Codecs()
			app.Logger.Info("conn=%v client=%v codecs=%v", conn, client, codecs)
		}
	}

	// Success
	return nil
}

func HasService(services []string, service string) bool {
	for _, value := range services {
		if value == service {
			return true
		}
	}
	return false
}

func Connect(app *gopi.AppInstance, service *gopi.RPCServiceRecord) error {

	// Start the connection process for a remote service
	clientpool := app.ModuleInstance("rpc/clientpool").(mutablehome.RPCClientPool)
	if client, err := clientpool.Connect(service, gopi.RPC_FLAG_NONE); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Add in discovery of "remotes" and "mihome" services
	if err := StartDiscovery(app, "remotes"); err != nil {
		done <- gopi.DONE
		return err
	}
	if err := StartDiscovery(app, "mihome"); err != nil {
		done <- gopi.DONE
		return err
	}

	app.Logger.Info("Press CTRL+C or send SIGTERM to terminate")
	app.WaitForSignal()

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Bootstrap

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/discovery", "rpc/clientpool", "rpc/client/mutablehome")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, DiscoveryLoop, EventLoop))
}
