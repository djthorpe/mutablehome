/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri_test

import (
	"net"
	"os"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	home "github.com/djthorpe/mutablehome"
)

func Test_Tradfri_000(t *testing.T) {
	t.Log("Test_Tradfri_000")
}

func Test_Tradfri_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_001, nil, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_001(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)
	t.Log(tradfri)
}

func Test_Tradfri_002(t *testing.T) {
	key, _ := os.LookupEnv("TRADFRI_KEY")
	args := []string{"-tradfri.id", "test3", "-tradfri.key", key}
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_002, args, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_002(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)

	if err := tradfri.Connect(gopi.RPCServiceRecord{
		Port: 5684,
		Addrs: []net.IP{
			net.ParseIP("192.168.86.56"),
		},
	}, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error(err)
	} else {
		t.Log(tradfri)
	}
}

func Test_Tradfri_003(t *testing.T) {
	key, _ := os.LookupEnv("TRADFRI_KEY")
	args := []string{"-tradfri.id", "test3", "-tradfri.key", key}
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_003, args, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_003(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)

	if err := tradfri.Connect(gopi.RPCServiceRecord{
		Port: 5684,
		Addrs: []net.IP{
			net.ParseIP("192.168.86.56"),
		},
	}, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error(err)
	} else if devices, err := tradfri.Devices(); err != nil {
		t.Error(err)
	} else if groups, err := tradfri.Groups(); err != nil {
		t.Error(err)
	} else if scenes, err := tradfri.Scenes(); err != nil {
		t.Error(err)
	} else {
		t.Log("devices=", devices)
		t.Log("groups=", groups)
		t.Log("scenes=", scenes)
	}
}

func Test_Tradfri_004(t *testing.T) {
	key, _ := os.LookupEnv("TRADFRI_KEY")
	args := []string{"-tradfri.id", "test3", "-tradfri.key", key}
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_004, args, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_004(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)

	if err := tradfri.Connect(gopi.RPCServiceRecord{
		Port: 5684,
		Addrs: []net.IP{
			net.ParseIP("192.168.86.56"),
		},
	}, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error(err)
	} else if devices, err := tradfri.Devices(); err != nil {
		t.Error(err)
	} else {
		for _, id := range devices {
			if device, err := tradfri.Device(id); err != nil {
				t.Error(err)
			} else {
				t.Log(device)
			}
		}
	}
}

func Test_Tradfri_005(t *testing.T) {
	key, _ := os.LookupEnv("TRADFRI_KEY")
	args := []string{"-tradfri.id", "test3", "-tradfri.key", key}
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_005, args, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_005(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)

	if err := tradfri.Connect(gopi.RPCServiceRecord{
		Port: 5684,
		Addrs: []net.IP{
			net.ParseIP("192.168.86.56"),
		},
	}, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error(err)
	} else if groups, err := tradfri.Groups(); err != nil {
		t.Error(err)
	} else {
		for _, id := range groups {
			if group, err := tradfri.Group(id); err != nil {
				t.Error(err)
			} else {
				t.Log(group)
			}
		}
	}
}

func Test_Tradfri_006(t *testing.T) {
	key, _ := os.LookupEnv("TRADFRI_KEY")
	args := []string{"-tradfri.id", "test3", "-tradfri.key", key}
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_006, args, "tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_006(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("tradfri").(home.Ikea)

	if err := tradfri.Connect(gopi.RPCServiceRecord{
		Port: 5684,
		Addrs: []net.IP{
			net.ParseIP("192.168.86.56"),
		},
	}, gopi.RPC_FLAG_INET_V4|gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error(err)
	} else if devices, err := tradfri.Devices(); err != nil {
		t.Error(err)
	} else {
		for _, id := range devices {
			if device, err := tradfri.Device(id); err != nil {
				t.Error(err)
			} else {
				for _, light := range device.Lights() {
					if light.Temperature() > 250 {
						if err := tradfri.Send(light.SetTemperature(250, 0)); err != nil {
							t.Error(err)
						}
					} else {
						if err := tradfri.Send(light.SetTemperature(450, 0)); err != nil {
							t.Error(err)
						}
					}
				}
			}
		}
	}
}
