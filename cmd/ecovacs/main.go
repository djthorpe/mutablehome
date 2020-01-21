package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	ecovacs := app.UnitInstance("mutablehome/ecovacs").(mutablehome.Ecovacs)

	if err := ecovacs.Authenticate(); err != nil {
		return err
	} else if devices, err := ecovacs.Devices(); err != nil {
		return err
	} else if len(devices) == 0 {
		return errors.New("No ecovacs devices found")
	} else {
		// For each device, connect and then clean for 10 seconds then return to charger
		for _, device := range devices {
			fmt.Println("Connect:", device)
			if err := ecovacs.Connect(device); err != nil {
				return err
			}
			if _, err := device.Clean(mutablehome.ECOVACS_CLEAN_AUTO, mutablehome.ECOVACS_SUCTION_STRONG); err != nil {
				return err
			} else {
				time.Sleep(10 * time.Second)
			}
			if _, err := device.Charge(); err != nil {
				return err
			} else {
				time.Sleep(2 * time.Second)
			}
		}

		// Wait for CTRL+C
		fmt.Println("Press CTRL+C to end")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}

	// Return success
	return nil
}
