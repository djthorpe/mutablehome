/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package client

import (

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	mutablehome "github.com/djthorpe/mutablehome"
	rpc "github.com/djthorpe/mutablehome/sys/rpc"

	// Protocol Buffer Implementations
	pb "github.com/djthorpe/remotes/protobuf/remotes"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// mutablelogic.Remotes
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/client/mutablehome",
		Type:     gopi.MODULE_TYPE_CLIENT,
		Requires: []string{"rpc/clientpool"},
		Run: func(app *gopi.AppInstance, _ gopi.Driver) error {
			// Register client functions with the clientpool
			clientpool := app.ModuleInstance("rpc/clientpool").(mutablehome.RPCClientPool)
			clientpool.RegisterClient("mutablelogic.Remotes", NewRemotesClient)
			// Return success
			return nil
		},
	})
}

func NewRemotesClient(conn mutablehome.RPCClientConn) (mutablehome.RPCClient, error) {
	return pb.NewRemotesClient(conn.(rpc.GRPCClientConn).Conn()), nil
}
