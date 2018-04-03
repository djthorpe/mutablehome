/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/mutablehome"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/mutablehome/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if devices := app.ModuleInstance("mutablehome/devices").(mutablehome.Devices); devices == nil {
		app.Logger.Error("Missing mutablehome/devices module")
		return gopi.ErrAppError
	} else {
		app.Logger.Info("devices=%v", devices)

		// get the device
		if device, err := devices.Device(0, mutablehome.DEVICE_TYPE_ENERGENIE_MONITOR, 0); err != nil {
			return err
		} else {
			fmt.Println(device)
		}
	}

	// Exit
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("mutablehome/devices")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
