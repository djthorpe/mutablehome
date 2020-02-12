/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	table := app.UnitInstance("mutablehome/dvb/table").(home.DVBTable)
	frontend := app.UnitInstance("mutablehome/dvb/frontend").(home.DVBFrontend)
	demux := app.UnitInstance("mutablehome/dvb/demux").(home.DVBDemux)

	prop := table.Properties()
	if len(prop) == 0 {
		return gopi.ErrNotFound
	}

	// Tune
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	fmt.Println("Tune", prop[0].Name(), "frequency=", prop[0].Frequency(), "Hz")
	if err := frontend.Tune(ctx, prop[0]); err != nil {
		return err
	}

	// Scan PAT
	if filter, err := demux.NewSectionFilter(0x0000); err != nil {
		return err
	} else {
		fmt.Println(filter)
	}

	fmt.Println("Wait for CTRL+C")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
