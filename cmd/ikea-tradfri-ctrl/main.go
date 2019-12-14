/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"
	mutablehome "github.com/djthorpe/mutablehome"
)

const (
	DISCOVERY_TIMEOUT_MS = 750
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {
	tradfri := app.ModuleInstance("mutablehome/ikea-tradfri").(mutablehome.IkeaTradfri)

	if len(services) == 0 {
		return fmt.Errorf("Gateway not found, use -addr to specify a gateway")
	} else if len(services) > 1 {
		return fmt.Errorf("Multiple gateways found, use -addr to specify a gateway")
	} else if err := tradfri.Connect(services[0], gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		return err
	} else if devices, err := tradfri.Devices(); err != nil {
		return err
	} else {
		for _, device := range devices {
			if device, err := tradfri.Device(device); err != nil {
				return err
			} else if device.Type() == mutablehome.IKEA_DEVICE_TYPE_LIGHT {
				fmt.Println(device, device.Lights())
			} else {
				fmt.Println(device)
			}
		}
		/*
			for _, group := range groups {
				if group, err := tradfri.Group(group); err != nil {
					return err
				} else {
					fmt.Println(group)
				}
			}
		*/
		/*
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
			defer cancel()
			if err := tradfri.ObserveDevice(ctx, 65537); err != nil {
				return err
			}
		*/

		// Success
		done <- gopi.DONE
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("mutablehome/ikea-tradfri", "rpc/clientpool", "discovery")

	// Gateway advertises as _coap._udp
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_TYPE, "coap")
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_FLAGS, gopi.RPC_FLAG_INET_UDP)

	// Run the command line tool
	os.Exit(rpc.Client(config, DISCOVERY_TIMEOUT_MS*time.Millisecond, Main))
}
