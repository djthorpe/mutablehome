package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/mutablehome/unit/tradfri"
)

/////////////////////////////////////////////////////////////////////

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "tradfri"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagString("tradfri.addr", "localhost", "Tradfri Gateway")
		os.Exit(app.Run())
	}
}
