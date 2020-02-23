package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mosquitto "github.com/djthorpe/mosquitto"
	mutablehome "github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	ecovacs := app.UnitInstance("ecovacs").(mutablehome.Ecovacs)
	mqtt := app.UnitInstance("mosquitto").(mosquitto.Client)

	if err := ecovacs.Authenticate(); err != nil {
		return err
	} else if devices, err := ecovacs.Devices(); err != nil {
		return err
	} else if len(devices) == 0 {
		return errors.New("No ecovacs devices found")
	} else if err := mqtt.Connect(); err != nil {
		return err
	} else {
		// For each device, connect
		for _, device := range devices {
			app.Log().Info("Connect", device)
			if err := ecovacs.Connect(device); err != nil {
				return err
			}
		}
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
