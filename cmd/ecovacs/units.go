package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
	"github.com/djthorpe/mutablehome"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/mutablehome/sys/ecovacs"
)

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "mutablehome.EcovacsEvent", Handler: EventHandler},
	}
	Header sync.Once
)

/////////////////////////////////////////////////////////////////////

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "mutablehome/ecovacs"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}

/////////////////////////////////////////////////////////////////////

func EventHandler(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	Header.Do(func() {
		fmt.Printf("%-15s %-40s\n", "EVENT", "VALUE")
		fmt.Printf("%-15s %-40s\n", strings.Repeat("-", 15), strings.Repeat("-", 40))
	})
	evt := evt_.(mutablehome.EcovacsEvent)
	type_ := strings.TrimPrefix(fmt.Sprint(evt.Type()), "ECOVACS_EVENT_")
	value_ := fmt.Sprint(evt.Value())
	fmt.Printf("%-15s %-40s\n", type_, value_)
}
