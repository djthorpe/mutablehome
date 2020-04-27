/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

func init() {
	// Tradfri Gateway Connector
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Tradfri{}.Name(),
		Requires: []string{"bus"},
		Config: func(app gopi.App) error {
			app.Flags().FlagString("tradfri.id", "", "Unique identifier")
			app.Flags().FlagString("tradfri.key", "", "Security code")
			app.Flags().FlagString("tradfri.path", ".tradfri", "State storage path")
			app.Flags().FlagDuration("tradfri.timeout", 0, "Connection timeout")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Tradfri{
				Id:      app.Flags().GetString("tradfri.id", gopi.FLAG_NS_DEFAULT),
				Key:     app.Flags().GetString("tradfri.key", gopi.FLAG_NS_DEFAULT),
				Path:    app.Flags().GetString("tradfri.path", gopi.FLAG_NS_DEFAULT),
				Timeout: app.Flags().GetDuration("tradfri.timeout", gopi.FLAG_NS_DEFAULT),
				Bus:     app.Bus(),
			}, app.Log().Clone(Tradfri{}.Name()))
		},
	})

	// Mutablehome node
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Node{}.Name(),
		Requires: []string{Tradfri{}.Name()},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Node{
				Tradfri: app.UnitInstance(Tradfri{}.Name()).(mutablehome.Ikea),
			}, app.Log().Clone(Node{}.Name()))
		},
	})
}
