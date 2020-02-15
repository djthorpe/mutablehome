/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/mutablehome/unit/dvb"
)

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "DVBSectionEvent", Handler: DVBSectionEventHandler},
	}
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "mutablehome/dvb/table", "mutablehome/dvb/frontend", "mutablehome/dvb/demux"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
