/*
  Tradfri: Interface to Ikea Tradfri

  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package tradfri_test

import (
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
	if app, err := app.NewTestTool(t, Main_Test_Tradfri_001, nil, "mutablehome/tradfri"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Tradfri_001(app gopi.App, t *testing.T) {
	tradfri := app.UnitInstance("mutablehome/tradfri").(home.Tradfri)
	t.Log(tradfri)
}
