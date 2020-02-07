package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "mutablehome.EcovacsEvent", Handler: PrintEvent},
		gopi.EventHandler{Name: "mutablehome.EcovacsEvent", Handler: LogEvent},
	}
	Header sync.Once
)

/////////////////////////////////////////////////////////////////////

func PrintEvent(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	Header.Do(func() {
		fmt.Printf("%-15s %-40s\n", "EVENT", "VALUE")
		fmt.Printf("%-15s %-40s\n", strings.Repeat("-", 15), strings.Repeat("-", 40))
	})
	evt := evt_.(mutablehome.EcovacsEvent)
	type_ := strings.TrimPrefix(fmt.Sprint(evt.Type()), "ECOVACS_EVENT_")
	value_ := fmt.Sprint(evt.Value())
	fmt.Printf("%-15s %-40s\n", type_, value_)
}

func LogEvent(ctx context.Context, app gopi.App, evt_ gopi.Event) {
	influxdb := app.UnitInstance("influxdb/v1").(mutablehome.InfluxDB)
	evt := evt_.(mutablehome.EcovacsEvent)

	rs := influxdb.NewResultSet(map[string]string{
		"device_id":       evt.Device().Address(),
		"device_nickname": evt.Device().Nickname(),
	})
	measurement := strings.TrimPrefix(fmt.Sprint(evt.Type()), "ECOVACS_EVENT_")
	switch evt.Type() {
	case mutablehome.ECOVACS_EVENT_BATTERYLEVEL:
		if err := rs.Add(measurement, map[string]interface{}{
			"value": evt.Value(),
		}); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_CLEANSTATE:
		measurement := "STATE"
		values := evt.Value().([]interface{})
		if err := rs.Add(measurement, map[string]interface{}{
			"mode":    values[0],
			"suction": values[1],
		}); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_CHARGESTATE:
		measurement := "STATE"
		if err := rs.Add(measurement, map[string]interface{}{
			"chargestate": evt.Value(),
		}); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_LIFESPAN:
		values := evt.Value().([]interface{})
		part := values[0].(mutablehome.EcovacsPart)
		if err := rs.Add(measurement, map[string]interface{}{
			string(part): float32(values[1].(uint)) * 100.0 / float32(values[2].(uint)),
		}); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_VERSION:
		if err := rs.Add(measurement, map[string]interface{}{
			"value": evt.Value(),
		}); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_ERROR:
		values := evt.Value().([]interface{})
		errNum := values[0].(uint)
		errDescription := values[1].(string)
		if err := rs.Add(measurement, map[string]interface{}{
			"err_num":         errNum,
			"err_description": errDescription,
		}); err != nil {
			app.Log().Error(err)
		}
	default:
		return
	}
	if err := influxdb.Write(rs); err != nil {
		app.Log().Error(err)
	}
}
