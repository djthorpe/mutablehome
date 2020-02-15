/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (

	// Frameworks
	"fmt"

	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type sectionevent struct {
	source  gopi.Unit
	filter  mutablehome.DVBFilter
	section mutablehome.DVBSection
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewSectionEvent(source gopi.Unit, filter mutablehome.DVBFilter, section mutablehome.DVBSection) mutablehome.DVBSectionEvent {
	return &sectionevent{source, filter, section}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (*sectionevent) Name() string {
	return "DVBSectionEvent"
}

func (*sectionevent) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *sectionevent) Source() gopi.Unit {
	return this.source
}

func (this *sectionevent) Value() interface{} {
	return this.Section()
}

func (this *sectionevent) Type() mutablehome.DVBTableType {
	return this.section.Type()
}

func (this *sectionevent) Filter() mutablehome.DVBFilter {
	return this.filter
}

func (this *sectionevent) Section() mutablehome.DVBSection {
	return this.section
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *sectionevent) String() string {
	return "<" + this.Name() +
		" filter=" + fmt.Sprint(this.filter) +
		" section=" + fmt.Sprint(this.section) +
		">"
}
