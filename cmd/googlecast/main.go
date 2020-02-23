package main

import (
	"context"
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

/////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to end")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
