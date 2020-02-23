/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Cast{}.Name(),
		Requires: []string{"gopi/mdns/discovery"},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Cast{
				Discovery: app.UnitInstance("gopi/mdns/discovery").(gopi.RPCServiceDiscovery),
			}, app.Log().Clone(Cast{}.Name()))
		},
	})
}
