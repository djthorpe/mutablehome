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
		gopi.EventHandler{Name: "ikea.Event", Handler: PrintEvent},
	}
	Header sync.Once
	Format = "%-10v %-15v %-5v %s\n"
)

/////////////////////////////////////////////////////////////////////

// Print event header
func PrintHeader() {
	Header.Do(func() {
		fmt.Printf(Format, "EVENT", "DEVICE", "ID", "VALUE")
		fmt.Printf(Format, strings.Repeat("-", 10), strings.Repeat("-", 15), strings.Repeat("-", 5), strings.Repeat("-", 25))
	})
}

// PrintEvent
func PrintEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	PrintHeader()
	evt := evt_.(mutablehome.IkeaEvent)
	device := evt.Device()
	etype := strings.TrimPrefix(fmt.Sprint(evt.Type()), "IKEA_EVENT_DEVICE_")
	dtype := strings.TrimPrefix(fmt.Sprint(device.Type()), "IKEA_DEVICE_TYPE_")
	value := fmt.Sprint(device)
	if device.Type() == mutablehome.IKEA_DEVICE_TYPE_LIGHT {
		value = fmt.Sprint(device.Lights())
	}
	fmt.Printf(Format, etype, dtype, device.Id(), value)
}
