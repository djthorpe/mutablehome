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
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {
	device := app.ModuleInstance("mutablehome/ikea")
	fmt.Println(services)
	fmt.Println(device)

	// Wait for CTRL+C
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	// Success
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("mutablehome/ikea", "rpc/clientpool", "discovery")

	config.AppFlags.SetParam(gopi.PARAM_SERVICE_TYPE, "coap")
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_TYPE, "coap")

	// Run the command line tool
	os.Exit(rpc.Client(config, 200*time.Millisecond, Main))
}
