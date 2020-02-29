package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
	tablewriter "github.com/olekukonko/tablewriter"
)

/////////////////////////////////////////////////////////////////////

func PrintDevices(app gopi.App) error {
	cast := app.UnitInstance("googlecast").(mutablehome.Cast)
	devices := cast.Devices()
	if len(devices) == 0 {
		return gopi.ErrNotFound.WithPrefix("No Chromecasts found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Model", "Service"})
	for _, device := range cast.Devices() {
		srv := "Idle"
		if device.State() > 0 {
			srv = device.Service()
		}
		table.Append([]string{
			device.Id(),
			device.Name(),
			device.Model(),
			srv,
		})
	}
	table.Render()

	return nil
}

func SetVolume(app gopi.App) error {
	cast := app.UnitInstance("googlecast").(mutablehome.Cast)
	devices := cast.Devices()
	for _, device := range devices {
		if err := device.LaunchAppWithId("5C292C3E"); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func Main(app gopi.App, args []string) error {
	timeout := app.Flags().GetDuration("timeout", gopi.FLAG_NS_DEFAULT)
	time.Sleep(timeout)

	if err := PrintDevices(app); err != nil {
		return err
	}

	// volume
	if err := SetVolume(app); err != nil {
		return err
	}

	if watch := app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT); watch {
		// Wait for CTRL+C
		fmt.Println("Press CTRL+C to end")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}

	// Return success
	return nil
}
