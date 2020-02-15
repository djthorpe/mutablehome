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

type SectionSDT struct {
	SectionHeader
	NetworkId uint16
	Services  []*Service
	CRC       []byte
}

type Service struct {
	Id           uint16
	EITSchedule  bool
	EITFollowing bool
	Status       uint8
	Scrambled    bool
	Descriptors  []*RowDescriptor
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewSDT(tid mutablehome.DVBTableType, r *TSReader) (section mutablehome.DVBSection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	this := &SectionSDT{
		SectionHeader: NewHeader(tid, r),
		NetworkId:     r.Uint16(),
	}

	// Increment
	r.Uint8() // Reserved for future use

	// Iterate through the SDT rows
	for r.IsEOF() == false {
		if r.Size() == 4 {
			this.CRC = r.Bytes(4)
		} else {
			service := &Service{
				Id: r.Uint16(),
			}
			schedule := r.Uint8()
			length := r.Uint16()
			service.EITSchedule = schedule&0x02 != 0x00
			service.EITFollowing = schedule&0x01 != 0x00
			service.Status = uint8(length & 0xE000 >> 13)
			service.Scrambled = length&0x100 != 0x00
			service.Descriptors = NewDescriptors(NewTSReader(r.Bytes(int(length & 0x0FFF))))
			this.Services = append(this.Services, service)
		}
	}

	section = this
	return
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionSDT) String() string {
	return "<SDT" +
		" original_network_id=" + fmt.Sprintf("0x%04X", this.NetworkId) +
		" services=" + fmt.Sprint(this.Services) +
		" crc=" + fmt.Sprintf("0x%06X", this.CRC) +
		" header=" + this.SectionHeader.String() +
		">"
}

func (this *Service) String() string {
	return "<Service" +
		" id=" + fmt.Sprint(this.Id) +
		" eit_schedule=" + fmt.Sprint(this.EITSchedule) +
		" eit_following=" + fmt.Sprint(this.EITFollowing) +
		" status=" + fmt.Sprint(this.Status) +
		" scrambled=" + fmt.Sprint(this.Scrambled) +
		" descriptors=" + fmt.Sprint(this.Descriptors) +
		">"
}
