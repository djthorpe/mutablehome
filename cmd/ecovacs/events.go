package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mosquitto"
	"github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

const (
	MQTT_TOPIC_BATTERYLEVEL = "%v/%v/batterylevel"
	MQTT_TOPIC_VERSION      = "%v/%v/version"
	MQTT_TOPIC_CLEANSTATE   = "%v/%v/cleanstate"
	MQTT_TOPIC_CHARGESTATE  = "%v/%v/chargestate"
	MQTT_TOPIC_LIFESPAN     = "%v/%v/lifespan"
)

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "ecovacs.Event", Handler: PrintEvovacsEvent},
		gopi.EventHandler{Name: "ecovacs.Event", Handler: PublishEvovacsEvent},
		gopi.EventHandler{Name: "mosquitto.Event", Handler: PrintMQTTEvent},
		gopi.EventHandler{Name: "mosquitto.Event", Handler: SubscribeMQTTEvent},
		gopi.EventHandler{Name: "mosquitto.Event", Handler: MessageMQTTEvent},
	}
	Header sync.Once
)

/////////////////////////////////////////////////////////////////////

type Command struct {
	Mode    string `json:"mode"`
	Suction string `json:"suction"`
}

/////////////////////////////////////////////////////////////////////

// Print event header
func PrintHeader() {
	Header.Do(func() {
		fmt.Printf("     %-15s %-40s\n", "EVENT", "VALUE")
		fmt.Printf("     %-15s %-40s\n", strings.Repeat("-", 15), strings.Repeat("-", 40))
	})
}

// Print MQTT event
func PrintMQTTEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	PrintHeader()
	evt := evt_.(mosquitto.Event)

	type_ := strings.TrimPrefix(fmt.Sprint(evt.Type()), "MOSQ_FLAG_EVENT_")
	value_ := fmt.Sprint(evt.Value())
	fmt.Printf("MQTT %-15s %-40s\n", type_, value_)
}

// Subscribe to the devices within the topic (ie, ecovacs/XXXXXXXX)
// so that messages can be received
func SubscribeMQTTEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	mqtt := app.UnitInstance("mosquitto").(mosquitto.Client)
	evt := evt_.(mosquitto.Event)
	topic := app.Flags().GetString("topic", gopi.FLAG_NS_DEFAULT)
	qos := mosquitto.OptQOS(app.Flags().GetInt("qos", gopi.FLAG_NS_DEFAULT))

	switch evt.Type() {
	case mosquitto.MOSQ_FLAG_EVENT_CONNECT:
		if evt.ReturnCode() != 0 {
			app.Log().Warn("Return code from CONNECT is", evt.ReturnCode())
		} else if _, err := mqtt.Subscribe(topic+"/+", qos); err != nil {
			app.Log().Error(err)
		}
	}
}

// Handle incoming messages and parse into a Command object
func MessageMQTTEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	evt := evt_.(mosquitto.Event)

	switch evt.Type() {
	case mosquitto.MOSQ_FLAG_EVENT_MESSAGE:
		var command Command
		if err := json.Unmarshal(evt.Data(), &command); err != nil {
			app.Log().Error(err)
		} else if err := Execute(app, evt.Topic(), command); err != nil {
			app.Log().Error(err)
		}
	}
}

// Print Ecovacs event
func PrintEvovacsEvent(_ context.Context, _ gopi.App, evt_ gopi.Event) {
	PrintHeader()
	evt := evt_.(mutablehome.EcovacsEvent)

	type_ := strings.TrimPrefix(fmt.Sprint(evt.Type()), "ECOVACS_EVENT_")
	value_ := fmt.Sprint(evt.Value())
	fmt.Printf("ECOV %-15s %-40s\n", type_, value_)
}

// Publish Ecovacs event to MQTT as influx line protocol format
func PublishEvovacsEvent(_ context.Context, app gopi.App, evt_ gopi.Event) {
	mqtt := app.UnitInstance("mosquitto").(mosquitto.Client)
	evt := evt_.(mutablehome.EcovacsEvent)
	topic := app.Flags().GetString("topic", gopi.FLAG_NS_DEFAULT)

	opts := []mosquitto.Opt{
		mosquitto.OptQOS(app.Flags().GetInt("qos", gopi.FLAG_NS_DEFAULT)),
		mosquitto.OptTag("device_id", evt.Device().Id()),
		mosquitto.OptTag("nickname", evt.Device().Nickname()),
	}

	switch evt.Type() {
	case mutablehome.ECOVACS_EVENT_VERSION:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_VERSION, topic, evt.Device().Id())
		fields := map[string]interface{}{
			"version": evt.Value(),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_BATTERYLEVEL:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_BATTERYLEVEL, topic, evt.Device().Id())
		fields := map[string]interface{}{
			"batterylevel": evt.Value(),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_CLEANSTATE:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_CLEANSTATE, topic, evt.Device().Id())
		values := evt.Value().([]interface{})
		fields := map[string]interface{}{
			"mode":    fmt.Sprint(values[0]),
			"suction": fmt.Sprint(values[1]),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_CHARGESTATE:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_CHARGESTATE, topic, evt.Device().Id())
		fields := map[string]interface{}{
			"chargestate": evt.Value(),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_ERROR:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_CLEANSTATE, topic, evt.Device().Id())
		values := evt.Value().([]interface{})
		fields := map[string]interface{}{
			"err_num":         values[0].(uint),
			"err_description": values[1].(string),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	case mutablehome.ECOVACS_EVENT_LIFESPAN:
		measurement := topic
		topic := fmt.Sprintf(MQTT_TOPIC_LIFESPAN, topic, evt.Device().Id())
		values := evt.Value().([]interface{})
		opts = append(opts, mosquitto.OptTag("part", fmt.Sprint(values[0].(mutablehome.EcovacsPart))))
		fields := map[string]interface{}{
			"lifespan": float32(values[1].(uint)) * 100.0 / float32(values[2].(uint)),
		}
		if _, err := mqtt.PublishInflux(topic, measurement, fields, opts...); err != nil {
			app.Log().Error(err)
		}
	}
}

func Execute(app gopi.App, topic string, command Command) error {
	ecovacs := app.UnitInstance("ecovacs").(mutablehome.Ecovacs)

	if device := DeviceForTopic(ecovacs, topic); device == nil {
		return gopi.ErrNotFound.WithPrefix(topic)
	} else if charge := ParseChargeCommand(command); charge {
		_, err := device.Charge()
		return err
	} else if mode, suction, err := ParseCleanCommand(command); err != nil {
		return err
	} else {
		_, err := device.Clean(mode, suction)
		return err
	}
}

func DeviceForTopic(ecovacs mutablehome.Ecovacs, topic string) mutablehome.EvovacsDevice {
	if parts := strings.Split(topic, "/"); len(parts) < 2 {
		return nil
	} else if devices, err := ecovacs.Devices(); err != nil {
		return nil
	} else {
		for _, device := range devices {
			if device.Id() == parts[len(parts)-1] {
				return device
			}
		}
	}

	// Return not found
	return nil
}

func ParseChargeCommand(command Command) bool {
	if charge := strings.ToUpper(command.Mode); charge == "CHARGE" {
		return true
	} else if charge == "" && command.Suction == "" {
		return true
	} else {
		return false
	}
}

func ParseCleanCommand(command Command) (mutablehome.EcovacsCleanMode, mutablehome.EcovacsCleanSuction, error) {
	mode := mutablehome.ECOVACS_CLEAN_STOP
	suction := mutablehome.ECOVACS_SUCTION_STRONG

	// Parse mode
	switch strings.ToUpper(command.Mode) {
	case "STOP":
		mode = mutablehome.ECOVACS_CLEAN_STOP
	case "AUTO":
		mode = mutablehome.ECOVACS_CLEAN_AUTO
	case "BORDER":
		mode = mutablehome.ECOVACS_CLEAN_BORDER
	case "SPOT":
		mode = mutablehome.ECOVACS_CLEAN_SPOT
	case "SINGLEROOM", "ROOM":
		mode = mutablehome.ECOVACS_CLEAN_ROOM
	default:
		return mode, suction, gopi.ErrBadParameter.WithPrefix(command.Mode)
	}

	// Parse suction
	switch strings.ToUpper(command.Suction) {
	case "STANDARD", "":
		suction = mutablehome.ECOVACS_SUCTION_STANDARD
	case "STRONG":
		suction = mutablehome.ECOVACS_SUCTION_STRONG
	default:
		return mode, suction, gopi.ErrBadParameter.WithPrefix(command.Suction)
	}

	// Success
	return mode, suction, nil
}
