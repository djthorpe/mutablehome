/*
  Mutablehome Automation: Rotel Amplifiers
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
	rotel := app.UnitInstance("mutablehome/rotel").(home.Rotel)

	fmt.Println(rotel)
	
	fmt.Println("Wait for CTRL+C")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
