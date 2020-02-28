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
