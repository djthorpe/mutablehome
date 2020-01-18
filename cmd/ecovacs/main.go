package main

import (
	"errors"
	"fmt"

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
			}
		}
	}

	// Return success
	return nil
}
