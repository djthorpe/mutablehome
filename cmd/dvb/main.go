/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	table := app.UnitInstance("mutablehome/dvb/table").(home.DVBTable)
	frontend := app.UnitInstance("mutablehome/dvb/frontend").(home.DVBFrontend)

	for _, prop := range table.Properties() {
		fmt.Println(prop)
	}

	fmt.Println(frontend)

	// Return success
	return nil
}
