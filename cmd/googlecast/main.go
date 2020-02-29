package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
	tablewriter "github.com/olekukonko/tablewriter"
)

/////////////////////////////////////////////////////////////////////

func DevicesWithId(uuid string, devices []mutablehome.CastDevice) []mutablehome.CastDevice {
	// Return all devices if no id specified
	if uuid == "" {
		return devices
	}

	// Make a map of the ids for lookup
	idmap := make(map[string]bool)
	for _, id := range strings.Split(uuid, ",") {
		key := strings.TrimSpace(id)
		if key != "" {
			idmap[key] = true
		}
	}

	// Filter devices
	devices_ := make([]mutablehome.CastDevice, 0, len(devices))
	for _, device := range devices {
		if _, exists := idmap[device.Id()]; exists {
			devices_ = append(devices_, device)
		} else if _, exists := idmap[device.Name()]; exists {
			devices_ = append(devices_, device)
		}
	}

	// Return filtered devices
	return devices_
}

func PrintDevices(devices []mutablehome.CastDevice) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Model", "Service"})
	for _, device := range devices {
		srv := "Idle"
		if device.State() > 0 {
			srv = device.Service()
		}
		table.Append([]string{
			device.Id(),
			device.Name(),
			device.Model(),
			srv,
		})
	}
	table.Render()
}

func ExecuteCommand(app gopi.App, devices []mutablehome.CastDevice, command string, args string) error {

	// Connect to chromecasts
	cast := app.UnitInstance("googlecast").(mutablehome.Cast)
	for _, device := range devices {
		if err := cast.Connect(device, gopi.RPC_FLAG_NONE); err != nil {
			return err
		}
	}

	// Run command for chromecasts
	for _, c := range Commands {
		if c.Name == strings.ToLower(command) {
			if args := c.Re.FindStringSubmatch(args); len(args) > 0 {
				if c.Func != nil {
					return c.Func(app, devices, args[1:])
				} else {
					return nil
				}
			} else {
				return fmt.Errorf("Syntax error: %s", c.Syntax)
			}
		}
	}

	// Return not found
	return gopi.ErrNotFound.WithPrefix(command)
}

func Main(app gopi.App, args []string) error {
	// Return devices
	cast := app.UnitInstance("googlecast").(mutablehome.Cast)
	timeout := app.Flags().GetDuration("timeout", gopi.FLAG_NS_DEFAULT)
	uuid := app.Flags().GetString("id", gopi.FLAG_NS_DEFAULT)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Filter devices and either display them with no arguments or
	// execute a command otherwise
	if devices, err := cast.Devices(ctx); err != nil {
		return err
	} else if devices = DevicesWithId(uuid, devices); len(devices) == 0 {
		return gopi.ErrNotFound.WithPrefix("No Chromecasts found")
	} else if len(args) == 0 {
		PrintDevices(devices)
	} else if err := ExecuteCommand(app, devices, args[0], strings.Join(args[1:], " ")); err != nil {
		return err
	}

	// Watch for events
	if watch := app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT); watch {
		// Wait for CTRL+C
		fmt.Println("Press CTRL+C to end")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}

	// Return success
	return nil
}
