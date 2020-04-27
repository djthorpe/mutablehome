/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	// Register NodeService
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     NodeService{}.Name(),
		Type:     gopi.UNIT_RPC_SERVICE,
		Requires: []string{"server"},
		Config: func(app gopi.App) error {
			// Set service type to _mutablehome._tcp
			app.Flags().SetString("service", gopi.FLAG_NS_SERVICE, "mutablehome")
			// Return success
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			// Create NodeService
			return gopi.New(NodeService{
				Server: app.UnitInstance("server").(gopi.RPCServer),
			}, app.Log().Clone(NodeService{}.Name()))
		},
	})
	// Register NodeClient
	gopi.UnitRegister(gopi.UnitConfig{
		Name: NodeClient{}.Name(),
		Type: gopi.UNIT_RPC_CLIENT,
		Stub: func(conn gopi.RPCClientConn) (gopi.RPCClientStub, error) {
			if unit, err := gopi.New(NodeClient{Conn: conn}, nil); err != nil {
				return nil, err
			} else {
				return unit.(gopi.RPCClientStub), nil
			}
		},
	})
}
