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
	"time"

	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SectionEIT struct {
	SectionHeader
	StreamId    uint16
	NetworkId   uint16
	LastSection uint8
	LastTable   uint8
	Events      []*Event
	CRC         []byte
}

type Event struct {
	Id          uint16
	Start       time.Time
	Duration    time.Duration
	Status      uint8
	Scrambled   bool
	Descriptors []*RowDescriptor
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewEIT(tid mutablehome.DVBTableType, r *TSReader) (section mutablehome.DVBSection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	this := &SectionEIT{
		SectionHeader: NewHeader(tid, r),
		StreamId:      r.Uint16(),
		NetworkId:     r.Uint16(),
		LastSection:   r.Uint8(),
		LastTable:     r.Uint8(),
		Events:        make([]*Event, 0),
	}

	for r.IsEOF() == false {
		if r.Size() == 4 {
			this.CRC = r.Bytes(4)
		} else {
			event := &Event{
				Id:       r.Uint16(),
				Start:    r.DateTime(),
				Duration: r.Duration(),
			}
			length := r.Uint16()
			event.Status = uint8(length & 0xE000 >> 13)
			event.Scrambled = length&0x100 != 0x00
			event.Descriptors = NewDescriptors(NewTSReader(r.Bytes(int(length & 0x0FFF))))
			this.Events = append(this.Events, event)
		}
	}

	// Return decoded section
	section = this
	return
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionEIT) String() string {
	return "<EIT" +
		" stream_id=" + fmt.Sprint(this.StreamId) +
		" network_id=" + fmt.Sprintf("0x%04X", this.NetworkId) +
		" last_section=" + fmt.Sprint(this.LastSection) +
		" last_table=" + fmt.Sprintf("0x%02X", this.LastTable) +
		" events=" + fmt.Sprint(this.Events) +
		" crc=" + fmt.Sprintf("0x%06X", this.CRC) +
		" header=" + this.SectionHeader.String() +
		">"
}

func (this *Event) String() string {
	return "<Event" +
		" id=" + fmt.Sprint(this.Id) +
		" start=" + fmt.Sprint(this.Start.Local()) +
		" duration=" + fmt.Sprint(this.Duration) +
		" status=" + fmt.Sprint(this.Status) +
		" scrambled=" + fmt.Sprint(this.Scrambled) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		">"
}
