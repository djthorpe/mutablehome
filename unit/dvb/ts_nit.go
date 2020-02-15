/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"fmt"

	// Frameworks
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SectionNIT struct {
	SectionHeader
	Descriptors []*RowDescriptor
	Streams     []*NITStream
	CRC         []byte
}

type NITStream struct {
	Pid         uint16
	NetworkId   uint16
	Descriptors []*RowDescriptor
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewNIT(tid mutablehome.DVBTableType, r *TSReader) (section mutablehome.DVBSection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	this := &SectionNIT{
		SectionHeader: NewHeader(tid, r),
		Streams:       make([]*NITStream, 0),
	}

	// Descriptors
	descLength := int(r.Uint16() & 0x0FFF)
	if descLength > 0 {
		this.Descriptors = NewDescriptors(NewTSReader(r.Bytes(descLength)))
	}

	r.Uint16() // Streams Length
	for r.IsEOF() == false {
		if r.Size() == 4 {
			this.CRC = r.Bytes(4)
		} else {
			stream := &NITStream{
				Pid:       r.Uint16(),
				NetworkId: r.Uint16(),
			}
			length := r.Uint16()
			stream.Descriptors = NewDescriptors(NewTSReader(r.Bytes(int(length & 0x0FFF))))
			this.Streams = append(this.Streams, stream)
		}
	}

	section = this
	return
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionNIT) String() string {
	return "<NIT" +
		" streams=" + fmt.Sprint(this.Streams) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		" crc=" + fmt.Sprintf("0x%06X", this.CRC) +
		" header=" + this.SectionHeader.String() +
		">"
}

func (this *NITStream) String() string {
	return "<NITStream" +
		" pid=" + fmt.Sprint(this.Pid) +
		" network_id=" + fmt.Sprintf("0x%04X", this.NetworkId) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		">"
}
