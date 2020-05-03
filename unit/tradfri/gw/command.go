/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
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

func NewLightState(device uint, state lightbulb) mutablehome.TradfriCommand {
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
		return "<tradfri.Command path=" + strconv.Quote(this.Path()) + " body=" + string(data) + ">"
	}
}
