/*
	Mutablehome Automation: Web Server
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"context"
	"net/url"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type HttpServer interface {

	// Serve a folder of content, returns the base URL
	// of the files being served
	ServeStatic(string) (*url.URL, error)

	// Stop serving with context
	Stop(context.Context) error

	// Implements gopi.Unit
	gopi.Unit
}
