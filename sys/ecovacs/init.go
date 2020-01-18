package ecovacs

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "mutablehome/ecovacs",
		Config: func(app gopi.App) error {
			app.Flags().FlagString("ecovacs.country", "au", "Ecovacs Country Code")
			app.Flags().FlagString("ecovacs.email", "", "Ecovacs Account Email")
			app.Flags().FlagString("ecovacs.password", "", "Ecovacs Account Password")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Ecovacs{
				Country:      app.Flags().GetString("ecovacs.country", gopi.FLAG_NS_DEFAULT),
				AccountId:    app.Flags().GetString("ecovacs.email", gopi.FLAG_NS_DEFAULT),
				PasswordHash: MD5String(app.Flags().GetString("ecovacs.password", gopi.FLAG_NS_DEFAULT)),
			}, app.Log().Clone("mutablehome/ecovacs"))
		},
	})
}
