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
	mutablehome "github.com/djthorpe/mutablehome"
	gateway "github.com/djthorpe/mutablehome/unit/tradfri/gw"
	node "github.com/djthorpe/mutablehome/unit/tradfri/node"
)

////////////////////////////////////////////////////////////////////////////////

type Node interface {
	mutablehome.Node

	// Connect to gateway
	Connect(gopi.RPCServiceRecord, gopi.RPCFlag) error
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	// Gateway Connector
	gopi.UnitRegister(gopi.UnitConfig{
		Name: gateway.Tradfri{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagString("tradfri.id", "", "Unique identifier")
			app.Flags().FlagString("tradfri.key", "", "Security code")
			app.Flags().FlagString("tradfri.state", ".tradfri", "State storage path")
			app.Flags().FlagDuration("tradfri.timeout", 0, "Connection timeout")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(gateway.Tradfri{
				Id:      app.Flags().GetString("tradfri.id", gopi.FLAG_NS_DEFAULT),
				Key:     app.Flags().GetString("tradfri.key", gopi.FLAG_NS_DEFAULT),
				Path:    app.Flags().GetString("tradfri.state", gopi.FLAG_NS_DEFAULT),
				Timeout: app.Flags().GetDuration("tradfri.timeout", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(gateway.Tradfri{}.Name()))
		},
	})

	// Node Connector
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     node.Node{}.Name(),
		Requires: []string{gateway.Tradfri{}.Name()},
		New: func(app gopi.App) (gopi.Unit, error) {
			tradfri := app.UnitInstance(gateway.Tradfri{}.Name()).(mutablehome.TradfriGateway)
			return gopi.New(node.Node{
				Gateway: tradfri,
			}, app.Log().Clone(node.Node{}.Name()))
		},
	})
}
