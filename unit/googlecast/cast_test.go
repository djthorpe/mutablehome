/*
	Mutablehome Automation: Googlecast
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package googlecast_test

import (
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	home "github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

func Test_Cast_000(t *testing.T) {
	t.Log("Test_Cast_000")
}

func Test_Cast_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Cast_001, nil, "googlecast"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Cast_001(app gopi.App, t *testing.T) {
	cast := app.UnitInstance("googlecast").(home.Cast)
	t.Log(cast)
	time.Sleep(time.Second * 5)
}
