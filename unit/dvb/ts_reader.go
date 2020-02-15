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
	"strings"
	"syscall"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type TSReader struct {
	buf []byte
	i   int
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewTSReader(buf []byte) *TSReader {
	this := new(TSReader)
	if buf == nil || len(buf) == 0 {
		return nil
	}
	this.buf = buf
	this.i = 0

	// Success
	return this
}

func NewTSReaderRead(fd uintptr, capacity int) (*TSReader, error) {
	this := new(TSReader)
	this.buf = make([]byte, capacity)
	this.i = 0

	// Perform the read
	if n, err := syscall.Read(int(fd), this.buf); err != nil {
		return nil, err
	} else {
		this.buf = this.buf[0:n]
	}
	// Success
	return this, nil
}

func (this *TSReader) SetLength(len int) {
	if len <= 0 || this.i >= len {
		panic(fmt.Errorf("SetLength: SectionReader overflow"))
	}
	this.buf = this.buf[0:len]
}

// IsEOF is true if index is equal the length of buffer
func (this *TSReader) IsEOF() bool {
	return this.i == len(this.buf)
}

// Size returns the remaining bytes to be read
func (this *TSReader) Size() int {
	return len(this.buf) - this.i
}

// Uint8 returns a byte
func (this *TSReader) Uint8() (v uint8) {
	if this.i >= len(this.buf) {
		panic(fmt.Errorf("Uint8: SectionReader overflow"))
	}
	v, this.i = this.buf[this.i], this.i+1
	return
}

// Uint16 returns a word
func (this *TSReader) Uint16() (v uint16) {
	if this.i+1 >= len(this.buf) {
		panic(fmt.Errorf("Uint16: SectionReader overflow"))
	}
	v, this.i = uint16(this.buf[this.i])<<8+uint16(this.buf[this.i+1]), this.i+2
	return
}

// Bytes returns a byte array
func (this *TSReader) Bytes(sz int) (v []byte) {
	if sz == 0 {
		return nil
	} else if this.i+sz > len(this.buf) {
		panic(fmt.Errorf("Bytes: SectionReader overflow"))
	}
	v, this.i = this.buf[this.i:this.i+sz], this.i+sz
	return
}

// DateTime returns a time.Date (UTC) (5 bytes)
func (this *TSReader) DateTime() time.Time {

	datetime := this.Bytes(5)

	// Check for zero time
	if datetime[0] == 0xFF && datetime[1] == 0xFF && datetime[2] == 0xFF && datetime[3] == 0xFF && datetime[4] == 0xFF {
		return time.Time{}
	}

	// Do crazy conversion
	mjd := float32(uint(datetime[0])<<8 + uint(datetime[1]))
	year := int((mjd - 15078.2) / 365.25)
	month := int((mjd - 14956.1 - float32(year)*365.25) / 30.6001)
	day := int(mjd) - 14956 - int(float32(year)*365.25) - int(float32(month)*30.6001)
	if month == 14 || month == 15 {
		year = year + 1
		month = month - 12
	}

	// Return time in UTC
	return time.Date(int(year+1900), time.Month(month-1), int(day), int(datetime[2]), int(datetime[3]), int(datetime[4]), 0, time.UTC)
}

// Duration returns a time.Duration (3 bytes)
func (this *TSReader) Duration() time.Duration {
	duration := this.Bytes(3)
	duration_ := time.Duration(decodeBCD(duration[0]))*time.Hour +
		time.Duration(decodeBCD(duration[1]))*time.Minute +
		time.Duration(decodeBCD(duration[2]))*time.Second
	return duration_
}

// String returns the stringified version
func (this *TSReader) String() string {
	return fmt.Sprintf("<TSReader i=%v len=%v %v>", this.i, len(this.buf), strings.ToUpper(hex.EncodeToString(this.buf)))
}

func decodeBCD(bcd byte) int {
	h := int(bcd) >> 4
	if h > 9 {
		return -1
	}
	l := int(bcd) & 0x0f
	if l > 9 {
		return -1
	}
	return h*10 + l
}
