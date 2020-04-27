/*
	Mutablehome Automation: Web Server
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package httpd_test

import (
	"net/http"
	"os"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	mutablehome "github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

func Test_Http_000(t *testing.T) {
	t.Log("Test_Http_000")
}

func Test_Http_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Http_001, nil, "httpd"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Http_001(app gopi.App, t *testing.T) {
	httpd := app.UnitInstance("httpd").(mutablehome.HttpServer)
	client := http.Client{}

	if folder, err := os.Getwd(); err != nil {
		t.Error(err)
	} else if url, err := httpd.ServeStatic(folder); err != nil {
		t.Error(err)
	} else if response, err := client.Get(url.String()); err != nil {
		t.Error(err)
	} else {
		t.Log(httpd)
		t.Log(url)
		t.Log(response)
	}
}
