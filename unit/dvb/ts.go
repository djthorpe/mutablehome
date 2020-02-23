/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"encoding/hex"
	"fmt"
	"strconv"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	mutablehome "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SectionHeader struct {
	TableId     mutablehome.DVBTableType
	ServiceId   uint16
	Version     uint8
	Current     bool
	Section     uint8
	LastSection uint8
}

type RowDescriptor struct {
	Tag    uint8
	Length uint8
	Data   []byte
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TS_PACKET_LENGTH   = 188
	TS_SECTION_BUFSIZE = 4096
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func TSRead(fd uintptr) (mutablehome.DVBSection, error) {
	// Read in the section, parse it and return it
	if reader, err := NewTSReaderRead(fd, TS_SECTION_BUFSIZE); err != nil {
		return nil, err
	} else {
		// Extract the table_id and the length
		tableId := mutablehome.DVBTableType(reader.Uint8())
		reader.SetLength(int(reader.Uint16()&0x0FFF) + 3)

		// Parse the table
		switch tableId {
		case mutablehome.DVB_TS_TABLE_PAT:
			return NewPAT(tableId, reader)
		case mutablehome.DVB_TS_TABLE_PMT:
			return NewPMT(tableId, reader)
		case mutablehome.DVB_TS_TABLE_SDT, mutablehome.DVB_TS_TABLE_SDT_OTHER:
			return NewSDT(tableId, reader)
		case mutablehome.DVB_TS_TABLE_NIT, mutablehome.DVB_TS_TABLE_NIT_OTHER:
			return NewNIT(tableId, reader)
		case mutablehome.DVB_TS_TABLE_EIT, mutablehome.DVB_TS_TABLE_EIT_OTHER:
			return NewEIT(tableId, reader)
		default:
			return nil, gopi.ErrUnexpectedResponse.WithPrefix(fmt.Sprint(tableId))
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// HEADER

func NewHeader(tableId mutablehome.DVBTableType, r *TSReader) SectionHeader {
	this := SectionHeader{
		TableId:     tableId,
		ServiceId:   r.Uint16(),
		Version:     r.Uint8(),
		Section:     r.Uint8(),
		LastSection: r.Uint8(),
	}
	this.Current = this.Version&0x01 != 0x00
	this.Version = (this.Version >> 1) & 0x1F
	return this
}

func (this *SectionHeader) Type() mutablehome.DVBTableType {
	return this.TableId
}

////////////////////////////////////////////////////////////////////////////////
// DESCRIPTORS

func NewDescriptors(r *TSReader) []*RowDescriptor {
	if r == nil {
		return nil
	}
	descriptors := make([]*RowDescriptor, 0)
	for r.IsEOF() == false {
		row := &RowDescriptor{
			Tag:    r.Uint8(),
			Length: r.Uint8(),
		}
		row.Data = r.Bytes(int(row.Length))
		descriptors = append(descriptors, row)
	}
	return descriptors
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionHeader) String() string {
	return "<Header" +
		" service_id=" + fmt.Sprintf("0x%04X", this.ServiceId) +
		" version=" + fmt.Sprint(this.Version) +
		" current=" + fmt.Sprint(this.Current) +
		" section=" + fmt.Sprint(this.Section) +
		" last_section=" + fmt.Sprint(this.LastSection) +
		">"
}

func (this *RowDescriptor) String() string {
	str := "<"
	switch this.Tag {
	case 0x40:
		str += "network_name_descriptor"
		str += " name=" + strconv.Quote(string(this.Data))
	case 0x41:
		str += "service_list_descriptor"
		str += " name=" + hex.EncodeToString(this.Data)
	case 0x48:
		str += "service_descriptor"
		r := NewTSReader(this.Data)
		str += " type=" + fmt.Sprintf("0x%02X", r.Uint8())
		provider := r.Bytes(int(r.Uint8()))
		name := r.Bytes(int(r.Uint8()))
		str += " provider=" + strconv.Quote(string(provider))
		str += " name=" + strconv.Quote(string(name))
	case 0x4D:
		str += "short_event_descriptor"
		r := NewTSReader(this.Data)
		str += " lang=" + fmt.Sprintf("0x%06X", r.Bytes(3))
		name := r.Bytes(int(r.Uint8()))
		description := r.Bytes(int(r.Uint8()))
		str += " name=" + strconv.Quote(string(name))
		str += " description=" + strconv.Quote(string(description))
	case 0x54:
		str += "content_descriptor"
		str += " data=" + hex.EncodeToString(this.Data)
	case 0x55:
		str += "parental_rating_descriptor"
		str += " data=" + hex.EncodeToString(this.Data)
	case 0x5A:
		str += "terrestrial_delivery_system_descriptor"
		str += " data=" + hex.EncodeToString(this.Data)
	case 0x5F:
		str += "private_data_specifier_descriptor"
		str += " data=" + hex.EncodeToString(this.Data)
	case 0x73:
		str += "default_authority_descriptor"
		str += " data=" + hex.EncodeToString(this.Data)
	default:
		str += "service" +
			" tag=" + fmt.Sprintf("0x%02X", this.Tag) +
			" length=" + fmt.Sprint(this.Length) +
			" data=" + hex.EncodeToString(this.Data)
	}
	return str + ">"
}
