/*
  Tradfri: Interface to Ikea Tradfri

  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"fmt"
	"strconv"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tradfri struct {
	Id      string
	Key     string
	Path    string
	Timeout time.Duration
}

type tradfri struct {
	log     gopi.Logger
	id      string
	key     string
	path    string
	timeout time.Duration

	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

const (
	CONN_TIMEOUT = 5 * time.Second
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Tradfri) Name() string { return "mutablehome/tradfri" }

func (config Tradfri) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(tradfri)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *tradfri) Init(config Tradfri) error {
	this.id = config.Id
	this.key = config.Key

	// Set timeout
	if config.Timeout == 0 {
		this.timeout = CONN_TIMEOUT
	} else {
		this.timeout = config.Timeout
	}

	/*
		// Create path if it doesn't exist, and read token file
		if path, err := this.createPath(config.Path); err != nil {
			return nil, err
		} else if err := this.token.Read(path); err != nil {
			return nil, err
		} else {
			this.path = path
		}
	*/

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *tradfri) String() string {
	return "<" + this.Log.Name() +
		" id=" + strconv.Quote(this.id) +
		" key=" + strconv.Quote(this.key) +
		" path=" + strconv.Quote(this.path) +
		" timeout=" + fmt.Sprint(this.timeout) +
		">"
}
