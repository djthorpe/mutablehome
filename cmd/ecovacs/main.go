package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

func Main(app gopi.App, args []string) error {
	ecovacs := app.UnitInstance("mutablehome/ecovacs").(mutablehome.Ecovacs)

	if err := ecovacs.Authenticate(); err != nil {
		return err
	} else if devices, err := ecovacs.Devices(); err != nil {
		return err
	} else if len(devices) == 0 {
		return errors.New("No ecovacs devices found")
	} else {
		for _, device := range devices {
			fmt.Println("Connect:", device)
			if err := device.Connect(); err != nil {
				return err
			} else if err := device.FetchBatteryLevel(); err != nil {
				return err
			}
		}
	}

	// Wait for CTRL+C
	fmt.Println("Press CRTL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
