/*
	Rotel RS232 Control
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Command uint16
	Power   uint16
	Volume  uint16
	Source  uint16
	Mute    uint16
	Bypass  uint16
	Tone    int16
	Balance int16
	Dimmer  uint16
	Speaker uint16
	Update  uint16
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Rotel interface {
	gopi.Driver
	gopi.Publisher

	// Information
	Model() string

	// Get and set state
	Get() RotelState
	Set(RotelState) error

	// Send Command
	Send(Command) error
}

type RotelEvent interface {
	gopi.Event

	Type() EventType
	State() RotelState
}

type RotelState struct {
	Power
	Volume
	Mute
	Source
	Freq string
	Bypass
	Treble Tone
	Bass   Tone
	Balance
	Speaker
	Dimmer
	Update
}

type RotelClient interface {
	gopi.RPCClient
	gopi.Publisher

	// Ping remote service
	Ping() error

	// Get and set state
	Get() (RotelState, error)
	Set(RotelState) error

	// Send command
	Send(Command) error

	// Stream state changes
	StreamEvents(ctx context.Context) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ROTEL_POWER_NONE Power = 0
	ROTEL_POWER_ON   Power = iota
	ROTEL_POWER_STANDBY
	ROTEL_POWER_MAX = ROTEL_POWER_STANDBY
)

const (
	ROTEL_MUTE_NONE Mute = 0
	ROTEL_MUTE_ON   Mute = iota
	ROTEL_MUTE_OFF
	ROTEL_MUTE_MAX = ROTEL_MUTE_OFF
)

const (
	ROTEL_BYPASS_NONE Bypass = 0
	ROTEL_BYPASS_ON   Bypass = iota
	ROTEL_BYPASS_OFF
	ROTEL_BYPASS_MAX = ROTEL_BYPASS_OFF
)

const (
	ROTEL_SOURCE_NONE Source = 0
	ROTEL_SOURCE_CD   Source = iota
	ROTEL_SOURCE_COAX1
	ROTEL_SOURCE_COAX2
	ROTEL_SOURCE_OPT1
	ROTEL_SOURCE_OPT2
	ROTEL_SOURCE_AUX1
	ROTEL_SOURCE_AUX2
	ROTEL_SOURCE_TUNER
	ROTEL_SOURCE_PHONO
	ROTEL_SOURCE_USB
	ROTEL_SOURCE_BLUETOOTH
	ROTEL_SOURCE_PC_USB
	ROTEL_SOURCE_OTHER
	ROTEL_SOURCE_MAX = ROTEL_SOURCE_OTHER
)

const (
	ROTEL_VOLUME_NONE Volume = 0
	ROTEL_VOLUME_MIN  Volume = 1
	ROTEL_VOLUME_MAX  Volume = 96
)

const (
	ROTEL_SPEAKER_NONE Speaker = 0
	ROTEL_SPEAKER_A    Speaker = iota
	ROTEL_SPEAKER_B
	ROTEL_SPEAKER_ALL
	ROTEL_SPEAKER_OFF
)

const (
	ROTEL_TONE_NONE Tone = 0
	ROTEL_TONE_MIN  Tone = -100
	ROTEL_TONE_MAX  Tone = 100
	ROTEL_TONE_OFF  Tone = ROTEL_TONE_MAX + 1
)

const (
	ROTEL_BALANCE_NONE      Balance = 0
	ROTEL_BALANCE_LEFT_MAX  Balance = -15
	ROTEL_BALANCE_RIGHT_MAX Balance = 15
	ROTEL_BALANCE_OFF       Balance = ROTEL_BALANCE_RIGHT_MAX + 1
)

const (
	ROTEL_DIMMER_NONE Dimmer = 0
	ROTEL_DIMMER_MIN  Dimmer = 1
	ROTEL_DIMMER_MAX  Dimmer = 9
	ROTEL_DIMMER_OFF  Dimmer = ROTEL_DIMMER_MAX + 1
)

const (
	ROTEL_UPDATE_NONE   Update = 0
	ROTEL_UPDATE_MANUAL Update = 1
	ROTEL_UPDATE_AUTO   Update = 2
	ROTEL_UPDATE_OTHER  Update = 3
)

const (
	ROTEL_COMMAND_NONE Command = 0
	ROTEL_COMMAND_PLAY Command = iota
	ROTEL_COMMAND_STOP
	ROTEL_COMMAND_PAUSE
	ROTEL_COMMAND_TRACK_NEXT
	ROTEL_COMMAND_TRACK_PREV
	ROTEL_COMMAND_MUTE_TOGGLE
	ROTEL_COMMAND_VOL_UP
	ROTEL_COMMAND_VOL_DOWN
	ROTEL_COMMAND_BASS_UP
	ROTEL_COMMAND_BASS_DOWN
	ROTEL_COMMAND_BASS_RESET
	ROTEL_COMMAND_TREBLE_UP
	ROTEL_COMMAND_TREBLE_DOWN
	ROTEL_COMMAND_TREBLE_RESET
	ROTEL_COMMAND_BALANCE_LEFT
	ROTEL_COMMAND_BALANCE_RIGHT
	ROTEL_COMMAND_BALANCE_RESET
	ROTEL_COMMAND_SPEAKER_A_TOGGLE
	ROTEL_COMMAND_SPEAKER_B_TOGGLE
	ROTEL_COMMAND_DIMMER_TOGGLE
	ROTEL_COMMAND_POWER_TOGGLE
	ROTEL_COMMAND_MAX = ROTEL_COMMAND_POWER_TOGGLE
)

const (
	ROTEL_EVENT_TYPE_NONE  EventType = 0
	ROTEL_EVENT_TYPE_POWER EventType = iota
	ROTEL_EVENT_TYPE_VOLUME
	ROTEL_EVENT_TYPE_SOURCE
	ROTEL_EVENT_TYPE_MUTE
	ROTEL_EVENT_TYPE_FREQ
	ROTEL_EVENT_TYPE_BYPASS
	ROTEL_EVENT_TYPE_BASS
	ROTEL_EVENT_TYPE_TREBLE
	ROTEL_EVENT_TYPE_BALANCE
	ROTEL_EVENT_TYPE_SPEAKER
	ROTEL_EVENT_TYPE_DIMMER
	ROTEL_EVENT_TYPE_UPDATE
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s RotelState) String() string {
	parts := make([]string, 0, 10)
	if s.Power != ROTEL_POWER_NONE {
		parts = append(parts, " power=", strings.TrimPrefix(fmt.Sprint(s.Power), "ROTEL_POWER_"))
	}
	if s.Volume != ROTEL_VOLUME_NONE {
		parts = append(parts, " volume=", strings.TrimPrefix(fmt.Sprint(s.Volume), "ROTEL_VOLUME_"))
	}
	if s.Mute != ROTEL_MUTE_NONE {
		parts = append(parts, " mute=", strings.TrimPrefix(fmt.Sprint(s.Mute), "ROTEL_MUTE_"))
	}
	if s.Source != ROTEL_SOURCE_NONE {
		parts = append(parts, " source=", strings.TrimPrefix(fmt.Sprint(s.Source), "ROTEL_SOURCE_"))
	}
	if s.Freq != "" {
		parts = append(parts, " freq=", strconv.Quote(s.Freq))
	}
	if s.Bypass != ROTEL_BYPASS_NONE {
		parts = append(parts, " bypass=", strings.TrimPrefix(fmt.Sprint(s.Bypass), "ROTEL_BYPASS_"))
	}
	if s.Bass != ROTEL_TONE_NONE {
		parts = append(parts, " bass=", strings.TrimPrefix(fmt.Sprint(s.Bass), "ROTEL_TONE_"))
	}
	if s.Treble != ROTEL_TONE_NONE {
		parts = append(parts, " treble=", strings.TrimPrefix(fmt.Sprint(s.Treble), "ROTEL_TONE_"))
	}
	if s.Balance != ROTEL_BALANCE_NONE {
		parts = append(parts, " balance=", strings.TrimPrefix(fmt.Sprint(s.Balance), "ROTEL_BALANCE_"))
	}
	if s.Speaker != ROTEL_SPEAKER_NONE {
		parts = append(parts, " speaker=", strings.TrimPrefix(fmt.Sprint(s.Speaker), "ROTEL_SPEAKER_"))
	}
	if s.Dimmer != ROTEL_DIMMER_NONE {
		parts = append(parts, " dimmer=", strings.TrimPrefix(fmt.Sprint(s.Dimmer), "ROTEL_DIMMER_"))
	}
	if s.Update != ROTEL_UPDATE_NONE {
		parts = append(parts, " update=", strings.TrimPrefix(fmt.Sprint(s.Update), "ROTEL_UPDATE_"))
	}
	return fmt.Sprintf("<rotel.State>{ %v }", strings.TrimSpace(strings.Join(parts, "")))
}

func (p Power) String() string {
	switch p {
	case ROTEL_POWER_NONE:
		return "ROTEL_POWER_NONE"
	case ROTEL_POWER_ON:
		return "ROTEL_POWER_ON"
	case ROTEL_POWER_STANDBY:
		return "ROTEL_POWER_STANDBY"
	default:
		return "[?? Invalid Power value]"
	}
}

func (s Speaker) String() string {
	switch s {
	case ROTEL_SPEAKER_NONE:
		return "ROTEL_SPEAKER_NONE"
	case ROTEL_SPEAKER_A:
		return "ROTEL_SPEAKER_A"
	case ROTEL_SPEAKER_B:
		return "ROTEL_SPEAKER_B"
	case ROTEL_SPEAKER_ALL:
		return "ROTEL_SPEAKER_ALL"
	case ROTEL_SPEAKER_OFF:
		return "ROTEL_SPEAKER_OFF"
	default:
		return "[?? Invalid Speaker value]"
	}
}

func (b Bypass) String() string {
	switch b {
	case ROTEL_BYPASS_NONE:
		return "ROTEL_BYPASS_NONE"
	case ROTEL_BYPASS_ON:
		return "ROTEL_BYPASS_ON"
	case ROTEL_BYPASS_OFF:
		return "ROTEL_BYPASS_OFF"
	default:
		return "[?? Invalid Bypass value ]"
	}

}

func (m Mute) String() string {
	switch m {
	case ROTEL_MUTE_NONE:
		return "ROTEL_MUTE_NONE"
	case ROTEL_MUTE_ON:
		return "ROTEL_MUTE_ON"
	case ROTEL_MUTE_OFF:
		return "ROTEL_MUTE_OFF"
	default:
		return "[?? Invalid Mute value]"
	}
}

func (v Volume) String() string {
	if v == ROTEL_VOLUME_NONE {
		return "ROTEL_VOLUME_NONE"
	} else if v == ROTEL_VOLUME_MAX {
		return "ROTEL_VOLUME_MAX"
	} else if v == ROTEL_VOLUME_MIN {
		return "ROTEL_VOLUME_MIN"
	} else if v > ROTEL_VOLUME_MIN && v < ROTEL_VOLUME_MAX {
		return fmt.Sprintf("ROTEL_VOLUME_%d", v)
	} else {
		return "[?? Invalid Volume value]"
	}
}

func (t Tone) String() string {
	switch {
	case t == ROTEL_TONE_NONE:
		return "ROTEL_TONE_NONE"
	case t == ROTEL_TONE_MAX:
		return "ROTEL_TONE_MAX"
	case t == ROTEL_TONE_MIN:
		return "ROTEL_TONE_MIN"
	case t == ROTEL_TONE_OFF:
		return "ROTEL_TONE_OFF"
	case t >= ROTEL_TONE_MIN && t < ROTEL_TONE_NONE:
		return fmt.Sprintf("ROTEL_TONE_MINUS_%d", -t)
	case t <= ROTEL_TONE_MAX && t > ROTEL_TONE_NONE:
		return fmt.Sprintf("ROTEL_TONE_PLUS_%d", t)
	default:
		return "[?? Invalid Tone value]"
	}
}

func (b Balance) String() string {
	switch {
	case b == ROTEL_BALANCE_NONE:
		return "ROTEL_BALANCE_NONE"
	case b == ROTEL_BALANCE_OFF:
		return "ROTEL_BALANCE_OFF"
	case b == ROTEL_BALANCE_LEFT_MAX:
		return "ROTEL_BALANCE_LEFT_MAX"
	case b == ROTEL_BALANCE_RIGHT_MAX:
		return "ROTEL_BALANCE_RIGHT_MAX"
	case b >= ROTEL_BALANCE_LEFT_MAX && b < ROTEL_BALANCE_NONE:
		return fmt.Sprintf("ROTEL_BALANCE_LEFT_%d", -b)
	case b <= ROTEL_BALANCE_RIGHT_MAX && b > ROTEL_BALANCE_NONE:
		return fmt.Sprintf("ROTEL_BALANCE_RIGHT_%d", b)
	default:
		return "[?? Invalid Balance value]"
	}
}

func (s Source) String() string {
	switch s {
	case ROTEL_SOURCE_NONE:
		return "ROTEL_SOURCE_NONE"
	case ROTEL_SOURCE_CD:
		return "ROTEL_SOURCE_CD"
	case ROTEL_SOURCE_COAX1:
		return "ROTEL_SOURCE_COAX1"
	case ROTEL_SOURCE_COAX2:
		return "ROTEL_SOURCE_COAX2"
	case ROTEL_SOURCE_OPT1:
		return "ROTEL_SOURCE_OPT1"
	case ROTEL_SOURCE_OPT2:
		return "ROTEL_SOURCE_OPT2"
	case ROTEL_SOURCE_AUX1:
		return "ROTEL_SOURCE_AUX1"
	case ROTEL_SOURCE_AUX2:
		return "ROTEL_SOURCE_AUX2"
	case ROTEL_SOURCE_TUNER:
		return "ROTEL_SOURCE_TUNER"
	case ROTEL_SOURCE_PHONO:
		return "ROTEL_SOURCE_PHONO"
	case ROTEL_SOURCE_USB:
		return "ROTEL_SOURCE_USB"
	case ROTEL_SOURCE_BLUETOOTH:
		return "ROTEL_SOURCE_BLUETOOTH"
	case ROTEL_SOURCE_PC_USB:
		return "ROTEL_SOURCE_PC_USB"
	case ROTEL_SOURCE_OTHER:
		return "ROTEL_SOURCE_OTHER"
	default:
		return "[?? Invalid Source value]"
	}
}

func (d Dimmer) String() string {
	switch {
	case d == ROTEL_DIMMER_NONE:
		return "ROTEL_DIMMER_NONE"
	case d == ROTEL_DIMMER_MAX:
		return "ROTEL_DIMMER_MAX"
	case d == ROTEL_DIMMER_MIN:
		return "ROTEL_DIMMER_MIN"
	case d == ROTEL_DIMMER_OFF:
		return "ROTEL_DIMMER_OFF"
	case d > ROTEL_DIMMER_MIN && d < ROTEL_DIMMER_MAX:
		return fmt.Sprintf("ROTEL_DIMMER_%d", d)
	default:
		return "[?? Invalid Dimmer value]"
	}
}

func (u Update) String() string {
	switch u {
	case ROTEL_UPDATE_NONE:
		return "ROTEL_UPDATE_NONE"
	case ROTEL_UPDATE_MANUAL:
		return "ROTEL_UPDATE_MANUAL"
	case ROTEL_UPDATE_AUTO:
		return "ROTEL_UPDATE_AUTO"
	default:
		return "[?? Invalid Update value]"
	}
}

/*
func (e EventType) String() string {
	switch e {
	case ROTEL_EVENT_TYPE_NONE:
		return "ROTEL_EVENT_TYPE_NONE"
	case ROTEL_EVENT_TYPE_POWER:
		return "ROTEL_EVENT_TYPE_POWER"
	case ROTEL_EVENT_TYPE_VOLUME:
		return "ROTEL_EVENT_TYPE_VOLUME"
	case ROTEL_EVENT_TYPE_SOURCE:
		return "ROTEL_EVENT_TYPE_SOURCE"
	case ROTEL_EVENT_TYPE_MUTE:
		return "ROTEL_EVENT_TYPE_MUTE"
	case ROTEL_EVENT_TYPE_FREQ:
		return "ROTEL_EVENT_TYPE_FREQ"
	case ROTEL_EVENT_TYPE_BYPASS:
		return "ROTEL_EVENT_TYPE_BYPASS"
	case ROTEL_EVENT_TYPE_BASS:
		return "ROTEL_EVENT_TYPE_BASS"
	case ROTEL_EVENT_TYPE_TREBLE:
		return "ROTEL_EVENT_TYPE_TREBLE"
	case ROTEL_EVENT_TYPE_BALANCE:
		return "ROTEL_EVENT_TYPE_BALANCE"
	case ROTEL_EVENT_TYPE_SPEAKER:
		return "ROTEL_EVENT_TYPE_SPEAKER"
	case ROTEL_EVENT_TYPE_DIMMER:
		return "ROTEL_EVENT_TYPE_DIMMER"
	case ROTEL_EVENT_TYPE_UPDATE:
		return "ROTEL_EVENT_TYPE_UPDATE"
	default:
		return "[?? Invalid EventType value]"
	}
}
*/

func (c Command) String() string {
	switch c {
	case ROTEL_COMMAND_NONE:
		return "ROTEL_COMMAND_NONE"
	case ROTEL_COMMAND_PLAY:
		return "ROTEL_COMMAND_PLAY"
	case ROTEL_COMMAND_STOP:
		return "ROTEL_COMMAND_STOP"
	case ROTEL_COMMAND_PAUSE:
		return "ROTEL_COMMAND_PAUSE"
	case ROTEL_COMMAND_TRACK_NEXT:
		return "ROTEL_COMMAND_TRACK_NEXT"
	case ROTEL_COMMAND_TRACK_PREV:
		return "ROTEL_COMMAND_TRACK_PREV"
	case ROTEL_COMMAND_MUTE_TOGGLE:
		return "ROTEL_COMMAND_MUTE_TOGGLE"
	case ROTEL_COMMAND_VOL_UP:
		return "ROTEL_COMMAND_VOL_UP"
	case ROTEL_COMMAND_VOL_DOWN:
		return "ROTEL_COMMAND_VOL_DOWN"
	case ROTEL_COMMAND_BASS_UP:
		return "ROTEL_COMMAND_BASS_UP"
	case ROTEL_COMMAND_BASS_DOWN:
		return "ROTEL_COMMAND_BASS_DOWN"
	case ROTEL_COMMAND_BASS_RESET:
		return "ROTEL_COMMAND_BASS_RESET"
	case ROTEL_COMMAND_TREBLE_UP:
		return "ROTEL_COMMAND_TREBLE_UP"
	case ROTEL_COMMAND_TREBLE_DOWN:
		return "ROTEL_COMMAND_TREBLE_DOWN"
	case ROTEL_COMMAND_TREBLE_RESET:
		return "ROTEL_COMMAND_TREBLE_RESET"
	case ROTEL_COMMAND_BALANCE_LEFT:
		return "ROTEL_COMMAND_BALANCE_LEFT"
	case ROTEL_COMMAND_BALANCE_RIGHT:
		return "ROTEL_COMMAND_BALANCE_RIGHT"
	case ROTEL_COMMAND_BALANCE_RESET:
		return "ROTEL_COMMAND_BALANCE_RESET"
	case ROTEL_COMMAND_SPEAKER_A_TOGGLE:
		return "ROTEL_COMMAND_SPEAKER_A_TOGGLE"
	case ROTEL_COMMAND_SPEAKER_B_TOGGLE:
		return "ROTEL_COMMAND_SPEAKER_B_TOGGLE"
	case ROTEL_COMMAND_DIMMER_TOGGLE:
		return "ROTEL_COMMAND_DIMMER_TOGGLE"
	case ROTEL_COMMAND_POWER_TOGGLE:
		return "ROTEL_COMMAND_POWER_TOGGLE"
	default:
		return "[?? Invalid Command value]"
	}
}
