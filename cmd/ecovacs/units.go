package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/mosquitto/unit/mosquitto"
	_ "github.com/djthorpe/mutablehome/unit/ecovacs"
)

/////////////////////////////////////////////////////////////////////

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "ecovacs", "mosquitto"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}
