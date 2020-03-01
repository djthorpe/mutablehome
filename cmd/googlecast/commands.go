package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

/////////////////////////////////////////////////////////////////////

type Command struct {
	Name   string
	Syntax string
	Re     *regexp.Regexp
	Func   func(gopi.App, []mutablehome.CastDevice, []string) error
}

var (
	Commands = []Command{
		Command{"volume", "volume <0-100>", regexp.MustCompile("^(\\d+)$"), Volume},
		Command{"app", "app <id>", regexp.MustCompile("^([A-Fa-f0-9]+)$"), LaunchApp},
		Command{"play", "play", regexp.MustCompile("^$"), Play},
		Command{"mute", "mute", regexp.MustCompile("^$"), Mute},
		Command{"unmute", "unmute", regexp.MustCompile("^$"), Unmute},
		Command{"pause", "pause", regexp.MustCompile("^$"), Pause},
		Command{"stop", "stop", regexp.MustCompile("^$"), Stop},
		Command{"load", "load <url>", regexp.MustCompile("^(http[s]?:.*)$"), Load},
	}
)

/////////////////////////////////////////////////////////////////////

func Load(_ gopi.App, devices []mutablehome.CastDevice, args []string) error {
	time.Sleep(2 * time.Second)
	fmt.Println("LOAD")
	for _, device := range devices {
		if err := device.LoadURL(args[0], true); err != nil {
			return err
		}
	}
	return nil
}

func Play(_ gopi.App, devices []mutablehome.CastDevice, _ []string) error {
	for _, device := range devices {
		if err := device.SetPlay(true); err != nil {
			return err
		}
	}
	return nil
}

func Pause(_ gopi.App, devices []mutablehome.CastDevice, _ []string) error {
	for _, device := range devices {
		if err := device.SetPause(true); err != nil {
			return err
		}
	}
	return nil
}

func Stop(_ gopi.App, devices []mutablehome.CastDevice, _ []string) error {
	for _, device := range devices {
		if err := device.SetPlay(false); err != nil {
			return err
		}
	}
	return nil
}

func Mute(_ gopi.App, devices []mutablehome.CastDevice, _ []string) error {
	for _, device := range devices {
		if err := device.SetMute(true); err != nil {
			return err
		}
	}
	return nil
}

func Unmute(_ gopi.App, devices []mutablehome.CastDevice, _ []string) error {
	for _, device := range devices {
		if err := device.SetMute(false); err != nil {
			return err
		}
	}
	return nil
}

func Volume(_ gopi.App, devices []mutablehome.CastDevice, args []string) error {
	if vol, err := strconv.ParseUint(args[0], 10, 32); err != nil {
		return err
	} else if vol > 100 {
		return gopi.ErrBadParameter.WithPrefix("volume")
	} else {
		level := float32(vol) / float32(100)
		for _, device := range devices {
			if err := device.SetVolume(level); err != nil {
				return err
			}
		}
	}
	return nil
}

func LaunchApp(_ gopi.App, devices []mutablehome.CastDevice, args []string) error {
	for _, device := range devices {
		if err := device.LaunchAppWithId(args[0]); err != nil {
			return err
		}
	}
	return nil
}
