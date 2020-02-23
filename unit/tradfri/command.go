/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Lightstate struct {
	Lightbulbs []lightbulb `json:"3311"`
}

type command struct {
	path []string
	body interface{}
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewLightState(device uint, state lightbulb) mutablehome.IkeaCommand {
	return &command{
		path: []string{PATH_DEVICES, fmt.Sprint(device)},
		body: Lightstate{[]lightbulb{state}},
	}
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PATH AND BODY

func (this *command) Path() string {
	return strings.Join(this.path, "/")
}

func (this *command) Body() (io.Reader, error) {
	if data, err := json.Marshal(this.body); err != nil {
		return nil, err
	} else if buf := bytes.NewBuffer(data); buf == nil {
		return nil, gopi.ErrInternalAppError
	} else {
		return buf, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *command) String() string {
	if data, err := json.Marshal(this.body); err != nil {
		return fmt.Sprint(err)
	} else {
		return "<IkeaCommand path=" + strconv.Quote(this.Path()) + " body=" + string(data) + ">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// VALIDATE

func boolToUint(value bool) uint {
	if value {
		return 1
	} else {
		return 0
	}
}

func durationToTransition(duration time.Duration) float32 {
	return float32(duration.Milliseconds()) / 100.0
}
