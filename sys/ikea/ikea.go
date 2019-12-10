/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package ikea

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	coap "github.com/go-ocf/go-coap"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Ikea struct {
}

type ikea struct {
	log  gopi.Logger
	coap *coap.Client
}

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Ikea) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("<ikea.Open>{ config=%+v }", config)

	this := new(ikea)
	this.log = logger
	this.coap = &coap.Client{}

	// Success
	return this, nil
}

func (this *ikea) Close() error {
	this.log.Debug("<ikea.Close>{ }")

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ikea) String() string {
	return fmt.Sprintf("<ikea>{ coap=%v }", this.coap)
}
