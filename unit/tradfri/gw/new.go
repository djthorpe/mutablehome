/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package gateway

import (
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tradfri struct {
	Id      string
	Key     string
	Path    string
	Timeout time.Duration
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (Tradfri) Name() string { return "mutablehome/tradfri/gateway" }

func (config Tradfri) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(gateway)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
