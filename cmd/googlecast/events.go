package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "cast.Event", Handler: PrintCastEvent},
		gopi.EventHandler{Name: "cast.Event", Handler: ConnectCastEvent},
	}
	Header sync.Once
)

/////////////////////////////////////////////////////////////////////

// Print event header
func PrintHeader() {
	Header.Do(func() {
		fmt.Printf("%-15s %-40s\n", "EVENT", "VALUE")
		fmt.Printf("%-15s %-40s\n", strings.Repeat("-", 15), strings.Repeat("-", 40))
	})
}

// Print cast event
func PrintCastEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	evt := evt_.(mutablehome.CastEvent)

	if watch := app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT); watch {
		PrintHeader()
		type_ := strings.TrimPrefix(fmt.Sprint(evt.Type()), "CAST_EVENT_")
		value_ := fmt.Sprint(evt.Value())
		fmt.Printf("%-15s %-40s\n", type_, value_)
	}
}

// Connect to chromecasts
func ConnectCastEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	evt := evt_.(mutablehome.CastEvent)
	cast := app.UnitInstance("googlecast").(mutablehome.Cast)
	watch := app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT)

	// We only connect in watch mode
	switch {
	case watch == false:
		break
	case evt.Type() != mutablehome.CAST_EVENT_ADDED:
		break
	default:
		if err := cast.Connect(evt.Device(), gopi.RPC_FLAG_NONE); err != nil {
			fmt.Printf("%-15s %-40v\n", "ERROR", err)
		}
	}
}
