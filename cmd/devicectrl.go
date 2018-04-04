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
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/mutablehome/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if device_db := app.ModuleInstance("mutablehome/devices").(mutablehome.Devices); device_db == nil {
		app.Logger.Error("Missing mutablehome/devices module")
		return gopi.ErrAppError
	} else {

		fmt.Println(mutablehome.DEVICE_TYPE_NONE, uint(mutablehome.DEVICE_TYPE_NONE))
		fmt.Println(mutablehome.DEVICE_TYPE_ENERGENIE_CONTROL, uint(mutablehome.DEVICE_TYPE_ENERGENIE_CONTROL))
		fmt.Println(mutablehome.DEVICE_TYPE_ENERGENIE_MONITOR, uint(mutablehome.DEVICE_TYPE_ENERGENIE_MONITOR))
		fmt.Println(mutablehome.DEVICE_TYPE_ANY, uint(mutablehome.DEVICE_TYPE_ANY))

		// Make a fake device
		if _, err := device_db.Device(0, mutablehome.DEVICE_TYPE_ENERGENIE_MONITOR, 0); err != nil {
			return err
		}

		// Print out a list of all devices
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Device", "Name", "Product", "Type", "Status"})
		devices := device_db.Devices(mutablehome.DEVICE_TYPE_ANY, mutablehome.PAIR_STATUS_ANY)
		for _, device := range devices {
			table.Append([]string{
				device.Hash(),
				device.Name,
				fmt.Sprint(device.ProductId),
				fmt.Sprint(device.Type),
				fmt.Sprint(device.PairStatus),
			})
		}
		table.Render()
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
