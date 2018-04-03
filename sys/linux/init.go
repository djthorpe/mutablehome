/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register mutablehome/devices module
	gopi.RegisterModule(gopi.Module{
		Name: "mutablehome/devices",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("mutablehome.root", "/var/local/mutablehome", "Folder for mutablehome data")
			config.AppFlags.FlagString("mutablehome.devices", "devices.json", "Filename for mutablehome devices")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			root, _ := app.AppFlags.GetString("mutablehome.root")
			filename, _ := app.AppFlags.GetString("mutablehome.devices")
			return gopi.Open(Devices{
				Root:     root,
				Filename: filename,
			}, app.Logger)
		},
	})
}
