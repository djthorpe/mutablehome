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

type SectionPAT struct {
	SectionHeader
	Programs []*Program
	CRC      []byte
}

type Program struct {
	Program uint16
	Pid     uint16
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewPAT(tid mutablehome.DVBTableType, r *TSReader) (section mutablehome.DVBSection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	this := &SectionPAT{
		SectionHeader: NewHeader(tid, r),
		Programs:      make([]*Program, 0),
	}

	// Iterate through the PAT rows
	for r.IsEOF() == false {
		if r.Size() == 4 {
			this.CRC = r.Bytes(4)
		} else {
			this.Programs = append(this.Programs, &Program{
				Program: r.Uint16(),
				Pid:     r.Uint16() & 0x1FFF,
			})
		}
	}
	section = this
	return
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionPAT) String() string {
	return "<PAT" +
		" programs=" + fmt.Sprint(this.Programs) +
		" crc=" + fmt.Sprintf("0x%06X", this.CRC) +
		" header=" + this.SectionHeader.String() +
		">"
}

func (this *Program) String() string {
	return "<Program" +
		" program=" + fmt.Sprint(this.Program) +
		" pid=" + fmt.Sprint(this.Pid) +
		">"
}
