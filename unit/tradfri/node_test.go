package tradfri_test

import (
	"context"
	"testing"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	tradfri "github.com/djthorpe/mutablehome/unit/tradfri"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

////////////////////////////////////////////////////////////////////////////////

func Test_Node_000(t *testing.T) {
	t.Log("Test_Node_000")
}

////////////////////////////////////////////////////////////////////////////////

func Test_Node_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Node_001, nil, "mutablehome/tradfri/node", "discovery"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Node_001(app gopi.App, t *testing.T) {
	node := app.UnitInstance("mutablehome/tradfri/node").(tradfri.Node)
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if services, err := discovery.Lookup(ctx, "_coap._udp"); err != nil {
		t.Error(err)
	} else if len(services) == 0 {
		t.Log("No gateway found, skipping tests")
	} else if err := node.Connect(services[0], 0); err != nil {
		t.Error(err)
	} else {
		t.Log(node)
		time.Sleep(2 * time.Second)
	}
}
