package main

import (
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"

	// Modules
	_ "github.com/djthorpe/gopi-rpc/sys/dns-sd"
	_ "github.com/djthorpe/gopi/sys/logger"

	// RPC
	_ "github.com/djthorpe/mutablehome/rpc/grpc/mutablehome"
)

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/mutablehome:service", "discovery")

	// Set subtype
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_SUBTYPE, "mutablehome")

	// Run the server and register all the services
	os.Exit(rpc.Server(config))
}
