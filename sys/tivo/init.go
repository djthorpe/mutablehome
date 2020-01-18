/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2019
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package tivo

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register TiVo
	gopi.RegisterModule(gopi.Module{
		Name: "tivo",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("tivo.mak", "", "Media Access Key")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			mak, _ := app.AppFlags.GetString("tivo.mak")
			return gopi.Open(TiVo{
				MediaAccessKey: mak,
			}, app.Logger)
		},
	})
}
