package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	uuid "github.com/google/uuid"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

///////////////////////////////////////////////////////////////////////////////

func WebServer(app *gopi.AppInstance, start chan<- struct{}, stop <-chan struct{}) error {
	start <- gopi.DONE

	port, _ := app.AppFlags.GetUint("port")
	client_id, _ := app.AppFlags.GetString("client_id")
	sslkey, _ := app.AppFlags.GetString("sslkey")
	sslcert, _ := app.AppFlags.GetString("sslcert")
	server := NewServer(port, client_id, app.Logger)
	errors := make(chan error)

	// Check SSL parameters
	if sslkey != "" || sslcert != "" {
		if sslkey == "" || sslcert == "" {
			return fmt.Errorf("Missing -sslkey or -sslcert flags")
		} else if stat, err := os.Stat(sslkey); os.IsNotExist(err) || stat.Mode().IsRegular() == false {
			return fmt.Errorf("-sslkey: Does not exist or not a regular file")
		} else if stat, err := os.Stat(sslcert); os.IsNotExist(err) || stat.Mode().IsRegular() == false {
			return fmt.Errorf("-sslcert: Does not exist or not a regular file")
		} else {
			server.ssl = []string{sslcert, sslkey}
		}
	}

	// Serve in background
	go func() {
		app.Logger.Info("Serving: %v", server.Addr)
		app.Logger.Info("Client ID: %v", client_id)
		if err := server.Serve(); err != nil {
			errors <- err
		}
	}()

	// Add a device
	server.AddDevice(&GoogleActionDevice{
		Id:   "coffee_maker_001",
		Type: GOOGLE_TYPE_COFFEE_MAKER,
		Traits: []string{
			GOOGLE_TRAIT_ONOFF,
		},
		Name: &GoogleActionDeviceName{
			Name: "Coffee Maker",
		},
		Room:            "Kitchen",
		WillReportState: false,
	})

	// Wait for error or stop signal
	for {
		select {
		case err := <-errors:
			app.SendSignal()
			return err
		case <-stop:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		}
	}
}

///////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Wait for signal
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()
	// Success
	return nil
}

func main() {
	// Create the configuration
	config := gopi.NewAppConfig()

	// Make a random default client id
	client_id := uuid.New().String()[:8]

	// Command line flags
	config.AppFlags.FlagUint("port", 9001, "Port for webserver")
	config.AppFlags.FlagString("sslcert", "", "SSL Certificate")
	config.AppFlags.FlagString("sslkey", "", "SSL Key")
	config.AppFlags.FlagString("client_id", client_id, "Application Client ID")

	// Run the server and register all the services
	os.Exit(gopi.CommandLineTool2(config, Main, WebServer))
}
