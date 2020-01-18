/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// An example RPC Client tool
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	rpc "github.com/djthorpe/gopi-rpc"

	// Modules
	_ "github.com/djthorpe/gopi-rpc/sys/dns-sd"
	_ "github.com/djthorpe/gopi-rpc/sys/grpc"
	_ "github.com/djthorpe/gopi-rpc/sys/rpcutil"
	_ "github.com/djthorpe/gopi-rpc/sys/tivo"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, services []gopi.RPCServiceRecord, done chan<- struct{}) error {

	if len(services) == 0 {
		return gopi.ErrNotFound
	} else if len(services) > 1 {
		names := make([]string, len(services))
		for i, record := range services {
			names[i] = strconv.Quote(record.Name())
		}
		return fmt.Errorf("More than one TiVo connected, choose with -addr <name>, where name is one of: %v", strings.Join(names, ","))
	} else if conn, err := app.ClientPool.Connect(services[0], gopi.RPC_FLAG_NONE); err != nil {
		return err
	} else if tivo := app.ModuleInstance("tivo").(rpc.TiVo); tivo == nil {
		return fmt.Errorf("Missing TiVo module")
	} else if session, err := tivo.NewSession(conn); err != nil {
		return err
	} else {
		fmt.Println(session)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("tivo", "discovery")

	// Set the service type
	config.AppFlags.SetParam(gopi.PARAM_SERVICE_TYPE, "tivo-mindrpc")

	// Run the command line tool
	os.Exit(rpc.Client(config, 700*time.Millisecond, Main))
}
