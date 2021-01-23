/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	// Frameworks

	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	mutablehome "github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi-rpc/v2/grpc/gaffer"
	_ "github.com/djthorpe/gopi-rpc/v2/unit/grpc"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
	_ "github.com/djthorpe/mutablehome/grpc/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// CONNECT

func ConnectAddrPort(app gopi.App, addr, port string) (gopi.RPCClientConn, error) {
	clientpool := app.UnitInstance("clientpool").(gopi.RPCClientPool)
	if port_, err := strconv.ParseInt(port, 10, 32); err != nil {
		return nil, gopi.ErrBadParameter.WithPrefix("-addr")
	} else if addr_ := net.ParseIP(addr); addr_ == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("-addr")
	} else {
		return clientpool.ConnectAddr(addr_, uint16(port_))
	}
}

func ConnectStub(app gopi.App, addr string) (mutablehome.NodeStub, error) {
	clientpool := app.UnitInstance("clientpool").(gopi.RPCClientPool)
	if host, port, err := net.SplitHostPort(addr); err != nil {
		return nil, err
	} else if conn, err := ConnectAddrPort(app, host, port); err != nil {
		return nil, err
	} else if stub, ok := clientpool.CreateStub("mutablehome.Node", conn).(mutablehome.NodeStub); ok == false {
		return nil, gopi.ErrInternalAppError.WithPrefix("CreateStub")
	} else {
		return stub, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// RUN COMMAND

func Run(app gopi.App, stub mutablehome.NodeStub, args []string) (bool, error) {
	fmt.Println("RUN", stub)

	// Success
	return false, nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN

func Main(app gopi.App, args []string) error {
	addr := app.Flags().GetString("addr", gopi.FLAG_NS_DEFAULT)
	if client, err := ConnectStub(app, addr); err != nil {
		return err
	} else if err := client.Ping(context.Background()); err != nil {
		return err
	} else if wait, err := Run(app, client, args); err != nil {
		return err
	} else if wait {
		// Wait until CTRL+C pressed
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "clientpool", "mutablehome.Node"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagString("addr", "", "Service address or name")

		// Run and exit
		os.Exit(app.Run())
	}
}
