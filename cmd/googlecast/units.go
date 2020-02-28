package main

import (
	"fmt"
	"os"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
	_ "github.com/djthorpe/mutablehome/unit/googlecast"
)

/////////////////////////////////////////////////////////////////////

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "googlecast"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagDuration("timeout", 500*time.Millisecond, "Discovery timeout")
		app.Flags().FlagBool("watch", false, "Watch for device changes")
		os.Exit(app.Run())
	}
}
