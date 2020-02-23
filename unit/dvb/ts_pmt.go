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

type SectionPMT struct {
	SectionHeader
	ClockPid    uint16
	Descriptors []*RowDescriptor
	Streams     []*PMTStream
	CRC         []byte
}

type PMTStream struct {
	Type        mutablehome.DVBStreamType
	Pid         uint16
	Length      uint16
	Descriptors []*RowDescriptor
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewPMT(tid mutablehome.DVBTableType, r *TSReader) (section mutablehome.DVBSection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	this := &SectionPMT{
		SectionHeader: NewHeader(tid, r),
		ClockPid:      r.Uint16() & 0x1FFF,
	}

	// Set other fields
	length := int(r.Uint16() & 0x3FF)
	if length > 0 {
		r := NewTSReader(r.Bytes(length))
		this.Descriptors = NewDescriptors(r)
	}

	// Iterate through the PAT streams
	this.Streams = make([]*PMTStream, 0)
	for r.IsEOF() == false {
		if r.Size() == 4 {
			this.CRC = r.Bytes(4)
		} else {
			stream := &PMTStream{
				Type:   mutablehome.DVBStreamType(r.Uint8()),
				Pid:    r.Uint16() & 0x1FFF,
				Length: r.Uint16() & 0x03FF,
			}
			stream.Descriptors = NewDescriptors(NewTSReader(r.Bytes(int(stream.Length))))
			this.Streams = append(this.Streams, stream)
		}
	}
	section = this
	return
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionPMT) String() string {
	return "<PMT" +
		" clock_pid=" + fmt.Sprintf("0x%04X", this.ClockPid) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		" streams=" + fmt.Sprint(this.Streams) +
		" crc=" + fmt.Sprintf("0x%06X", this.CRC) +
		" header=" + this.SectionHeader.String() +
		">"
}

func (this *PMTStream) String() string {
	return "<Stream" +
		" type=" + fmt.Sprint(this.Type) +
		" pid=" + fmt.Sprint(this.Pid) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		">"
}
