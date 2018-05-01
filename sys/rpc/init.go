/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register rpc/clientpool module
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/clientpool",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagBool("rpc.insecure", true, "Allow insecure SSL connections")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			insecure, _ := app.AppFlags.GetBool("rpc.insecure")
			return gopi.Open(ClientPool{
				SkipVerify: insecure,
			}, app.Logger)
		},
	})
}
