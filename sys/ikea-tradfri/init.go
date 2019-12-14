/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name: "mutablehome/ikea-tradfri",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("tradfri.id", "", "Unique identifier")
			config.AppFlags.FlagString("tradfri.key", "", "Security code")
			config.AppFlags.FlagString("tradfri.path", ".tradfri", "State storage path")
			config.AppFlags.FlagDuration("tradfri.timeout", 0, "Connection timeout")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			id, _ := app.AppFlags.GetString("tradfri.id")
			key, _ := app.AppFlags.GetString("tradfri.key")
			path, _ := app.AppFlags.GetString("tradfri.path")
			timeout, _ := app.AppFlags.GetDuration("tradfri.timeout")
			return gopi.Open(Tradfri{
				Id:      id,
				Key:     key,
				Path:    path,
				Timeout: timeout,
			}, app.Logger)
		},
	})
}
