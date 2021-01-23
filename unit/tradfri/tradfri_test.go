package tradfri_test

import (
	"context"
	"testing"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	mutablehome "github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

////////////////////////////////////////////////////////////////////////////////

func Test_Tradfri_000(t *testing.T) {
	t.Log("Test_Tradfri_000")
}

////////////////////////////////////////////////////////////////////////////////

func Test_Tradfri_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_001, nil, "mutablehome/tradfri/gateway"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_001(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("mutablehome/tradfri/gateway").(mutablehome.TradfriGateway)
	t.Log(tradfri)
}

////////////////////////////////////////////////////////////////////////////////

func Test_Tradfri_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_002, nil, "mutablehome/tradfri/gateway", "discovery"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_002(app gopi.App, t *testing.T) {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	tradfri := app.UnitInstance("mutablehome/tradfri/gateway").(mutablehome.TradfriGateway)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if services, err := discovery.Lookup(ctx, "_coap._udp"); err != nil {
		t.Error(err)
	} else if len(services) == 0 {
		t.Log("No gateway found, skipping tests")
	} else if err := tradfri.Connect(services[0], 0); err != nil {
		t.Error(err)
	} else {
		t.Log(tradfri)
	}
}

////////////////////////////////////////////////////////////////////////////////
func Test_Tradfri_003(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_003, nil, "mutablehome/tradfri/gateway", "discovery"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_003(app gopi.App, t *testing.T) {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	tradfri := app.UnitInstance("mutablehome/tradfri/gateway").(mutablehome.TradfriGateway)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if services, err := discovery.Lookup(ctx, "_coap._udp"); err != nil {
		t.Error(err)
	} else if len(services) == 0 {
		t.Log("No gateway found, skipping tests")
	} else if err := tradfri.Connect(services[0], 0); err != nil {
		t.Error(err)
	} else if devices, err := tradfri.Devices(); err != nil {
		t.Error(err)
	} else if groups, err := tradfri.Groups(); err != nil {
		t.Error(err)
	} else if scenes, err := tradfri.Scenes(); err != nil {
		t.Error(err)
	} else {
		t.Log("Devices=", devices)
		t.Log("Groups=", groups)
		t.Log("Scenes=", scenes)
	}
}
