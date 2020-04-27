package main

import (
	"context"
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	httpd := app.UnitInstance("httpd").(mutablehome.HttpServer)
	folder := ""
	if len(args) == 0 {
		if wd, err := os.Getwd(); err != nil {
			return err
		} else {
			folder = wd
		}
	} else if len(args) == 1 {
		folder = args[0]
	} else {
		return gopi.ErrHelp
	}

	// Start web server
	if url, err := httpd.ServeStatic(folder); err != nil {
		return err
	} else {
		fmt.Println("Serving on", url)
	}

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
