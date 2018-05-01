/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"
	// Frameworks
	gopi "github.com/djthorpe/gopi"
)

type clientevent struct {
	source  gopi.Driver
	type_   gopi.RPCEventType
	service *gopi.RPCServiceRecord
}

func (this *clientevent) Source() gopi.Driver {
	return this.source
}

func (this *clientevent) Name() string {
	return "RPCEvent"
}

func (this *clientevent) Type() gopi.RPCEventType {
	return this.type_
}

func (this *clientevent) ServiceRecord() *gopi.RPCServiceRecord {
	return this.service
}

func (this *clientevent) String() string {
	return fmt.Sprintf("<rpc.Event>{ type=%v source=%v }", this.Type(), this.Source())
}
