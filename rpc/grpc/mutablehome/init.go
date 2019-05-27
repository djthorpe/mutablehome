/*
	Mutablehome: Home Automation in Golang
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    https://github.com/djthorpe/mutablehome/
	For Licensing and Usage information, please see LICENSE
*/

package mutablehome

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register server
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/mutablehome:service",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(MutableHome{
				Server: app.ModuleInstance("rpc/server").(gopi.RPCServer),
			}, app.Logger)
		},
	})
}
