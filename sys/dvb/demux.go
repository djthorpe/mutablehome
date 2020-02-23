// +build linux

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
	"os"
	"time"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
	#include <sys/ioctl.h>
	#include <linux/dvb/dmx.h>
	static int _DMX_START() { return DMX_START; }
	static int _DMX_STOP() { return DMX_STOP; }
	static int _DMX_SET_FILTER() { return DMX_SET_FILTER; }
	static int _DMX_SET_PES_FILTER() { return DMX_SET_PES_FILTER; }
	static int _DMX_SET_BUFFER_SIZE() { return DMX_SET_BUFFER_SIZE; }
	static int _DMX_GET_PES_PIDS() { return DMX_GET_PES_PIDS; }
	static int _DMX_GET_STC() { return DMX_GET_STC; }
	static int _DMX_ADD_PID() { return DMX_ADD_PID; }
	static int _DMX_REMOVE_PID() { return DMX_REMOVE_PID; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DMXInput      uint32
	DMXOutput     uint32
	DMXStreamType uint32
	DMXFlags      uint32

	// DMXPattern specifies a section filter
	DMXPattern struct {
		Filter [16]byte
		Mask   [16]byte
		Mode   [16]byte
	}

	DMXStreamFilter struct {
		Pid   uint16
		In    DMXInput
		Out   DMXOutput
		Type  DMXStreamType
		Flags DMXFlags
	}

	DMXSectionFilter struct {
		Pid     uint16
		Pattern DMXPattern
		Timeout uint32 // Seconds, or zero for no timeout
		Flags   DMXFlags
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_DMX_IN_FRONTEND DMXInput = iota
	DVB_DMX_IN_DVR
)

const (
	DVB_DMX_OUT_DECODER DMXOutput = iota
	DVB_DMX_OUT_TAP
	DVB_DMX_OUT_TS_TAP
	DVB_DMX_OUT_TSDEMUX_TAP
)

const (
	DVB_DMX_PES_AUDIO0 DMXStreamType = iota
	DVB_DMX_PES_VIDEO0
	DVB_DMX_PES_TELETEXT0
	DVB_DMX_PES_SUBTITLE0
	DVB_DMX_PES_PCR0
	DVB_DMX_PES_AUDIO1
	DVB_DMX_PES_VIDEO1
	DVB_DMX_PES_TELETEXT1
	DVB_DMX_PES_SUBTITLE1
	DVB_DMX_PES_PCR1
	DVB_DMX_PES_AUDIO2
	DVB_DMX_PES_VIDEO2
	DVB_DMX_PES_TELETEXT2
	DVB_DMX_PES_SUBTITLE2
	DVB_DMX_PES_PCR2
	DVB_DMX_PES_AUDIO3
	DVB_DMX_PES_VIDEO3
	DVB_DMX_PES_TELETEXT3
	DVB_DMX_PES_SUBTITLE3
	DVB_DMX_PES_PCR3
	DVB_DMX_PES_OTHER
)

const (
	DVB_DMX_FLAG_NONE            DMXFlags = 0
	DVB_DMX_FLAG_CHECK_CRC       DMXFlags = 1
	DVB_DMX_FLAG_ONESHOT         DMXFlags = 2
	DVB_DMX_FLAG_IMMEDIATE_START DMXFlags = 4
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	DVB_DMX_START           = uintptr(C._DMX_START())
	DVB_DMX_STOP            = uintptr(C._DMX_STOP())
	DVB_DMX_SET_FILTER      = uintptr(C._DMX_SET_FILTER())
	DVB_DMX_SET_PES_FILTER  = uintptr(C._DMX_SET_PES_FILTER())
	DVB_DMX_SET_BUFFER_SIZE = uintptr(C._DMX_SET_BUFFER_SIZE())
	DVB_DMX_GET_PES_PIDS    = uintptr(C._DMX_GET_PES_PIDS())
	DVB_DMX_GET_STC         = uintptr(C._DMX_GET_STC())
	DVB_DMX_ADD_PID         = uintptr(C._DMX_ADD_PID())
	DVB_DMX_REMOVE_PID      = uintptr(C._DMX_REMOVE_PID())
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: DEMUX

func DVB_DMXPath(bus, demux uint) string {
	return fmt.Sprintf("%v%v/demux%v", DVB_PATH_WILDCARD, bus, demux)
}

func DVB_DVRPath(bus, demux uint) string {
	return fmt.Sprintf("%v%v/dvr%v", DVB_PATH_WILDCARD, bus, demux)
}

func DVB_DMXOpen(bus, demux uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_DMXPath(bus, demux), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func DVB_DVROpen(bus, demux uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_DVRPath(bus, demux), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func DVB_DMXStart(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_START, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("DVB_DMX_START", err)
	} else {
		return nil
	}
}

func DVB_DMXStop(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_STOP, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("DVB_DMX_STOP", err)
	} else {
		return nil
	}
}

func DVB_DMXSetBufferSize(fd uintptr, size uint32) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_BUFFER_SIZE, unsafe.Pointer(uintptr(size))); err != 0 {
		return os.NewSyscallError("DVB_DMX_SET_BUFFER_SIZE", err)
	} else {
		return nil
	}
}

func DVB_DMXSetSectionFilter(fd uintptr, filter DMXSectionFilter) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_FILTER, unsafe.Pointer(&filter)); err != 0 {
		return os.NewSyscallError("DVB_DMX_SET_FILTER", err)
	} else {
		return nil
	}
}

func DVB_DMXSetStreamFilter(fd uintptr, filter DMXStreamFilter) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_PES_FILTER, unsafe.Pointer(&filter)); err != 0 {
		return os.NewSyscallError("DVB_DMX_SET_PES_FILTER", err)
	} else {
		return nil
	}
}

func DVB_DMXAddPid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DVB_DMX_ADD_PID, unsafe.Pointer(&pid)); err != 0 {
		return os.NewSyscallError("DVB_DMX_ADD_PID", err)
	} else {
		return nil
	}
}

func DVB_DMXRemovePid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DVB_DMX_REMOVE_PID, unsafe.Pointer(&pid)); err != 0 {
		return os.NewSyscallError("DVB_DMX_REMOVE_PID", err)
	} else {
		return nil
	}
}

func DVB_DMXGetStreamPids(fd uintptr) (map[DMXStreamType]uint16, error) {
	var pids [5]uint16
	if err := dvb_ioctl(fd, DVB_DMX_GET_PES_PIDS, unsafe.Pointer(&pids)); err != 0 {
		return nil, os.NewSyscallError("DVB_DMX_GET_PES_PIDS", err)
	}
	pidmap := make(map[DMXStreamType]uint16)
	for stream, pid := range pids {
		if pid != uint16(0xFFFF) {
			pidmap[DMXStreamType(stream)] = pid
		}
	}
	return pidmap, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f DMXSectionFilter) String() string {
	return "<DVBSectionFilter" +
		fmt.Sprintf(" pid=0x%04X", f.Pid) +
		fmt.Sprintf(" timeout=%v", time.Second*time.Duration(f.Timeout)) +
		fmt.Sprintf(" filter=%s", hex.EncodeToString(f.Pattern.Filter[:])) +
		fmt.Sprintf(" mask=%s", hex.EncodeToString(f.Pattern.Mask[:])) +
		fmt.Sprintf(" flags=%v", f.Flags) +
		">"
}

func (f DMXStreamFilter) String() string {
	return "<DVBStreamFilter" +
		fmt.Sprintf(" pid=0x%04X", f.Pid) +
		fmt.Sprintf(" in=%v", f.In) +
		fmt.Sprintf(" out=%v", f.Out) +
		fmt.Sprintf(" type=%s", f.Type) +
		fmt.Sprintf(" flags=%v", f.Flags) +
		">"
}

func (s DMXStreamType) String() string {
	switch s {
	case DVB_DMX_PES_AUDIO0:
		return "DMX_PES_AUDIO0"
	case DVB_DMX_PES_VIDEO0:
		return "DVB_DMX_PES_VIDEO0"
	case DVB_DMX_PES_TELETEXT0:
		return "DVB_DMX_PES_TELETEXT0"
	case DVB_DMX_PES_SUBTITLE0:
		return "DVB_DMX_PES_SUBTITLE0"
	case DVB_DMX_PES_PCR0:
		return "DVB_DMX_PES_PCR0"
	case DVB_DMX_PES_AUDIO1:
		return "DVB_DMX_PES_AUDIO1"
	case DVB_DMX_PES_VIDEO1:
		return "DVB_DMX_PES_VIDEO1"
	case DVB_DMX_PES_TELETEXT1:
		return "DVB_DMX_PES_TELETEXT1"
	case DVB_DMX_PES_SUBTITLE1:
		return "DVB_DMX_PES_SUBTITLE1"
	case DVB_DMX_PES_PCR1:
		return "DVB_DMX_PES_PCR1"
	case DVB_DMX_PES_AUDIO2:
		return "DVB_DMX_PES_AUDIO2"
	case DVB_DMX_PES_VIDEO2:
		return "DVB_DMX_PES_VIDEO2"
	case DVB_DMX_PES_TELETEXT2:
		return "DVB_DMX_PES_TELETEXT2"
	case DVB_DMX_PES_SUBTITLE2:
		return "DVB_DMX_PES_SUBTITLE2"
	case DVB_DMX_PES_PCR2:
		return "DVB_DMX_PES_PCR2"
	case DVB_DMX_PES_AUDIO3:
		return "DVB_DMX_PES_AUDIO3"
	case DVB_DMX_PES_VIDEO3:
		return "DVB_DMX_PES_VIDEO3"
	case DVB_DMX_PES_TELETEXT3:
		return "DVB_DMX_PES_TELETEXT3"
	case DVB_DMX_PES_SUBTITLE3:
		return "DVB_DMX_PES_SUBTITLE3"
	case DVB_DMX_PES_PCR3:
		return "DVB_DMX_PES_PCR3"
	case DVB_DMX_PES_OTHER:
		return "DVB_DMX_PES_OTHER"
	default:
		return fmt.Sprintf("[?? Invalid DMXStreamType %02X]", uint8(s))
	}
}

func (f DMXFlags) String() string {
	switch f {
	case DVB_DMX_FLAG_NONE:
		return "DVB_DMX_FLAG_NONE"
	case DVB_DMX_FLAG_CHECK_CRC:
		return "DVB_DMX_FLAG_CHECK_CRC"
	case DVB_DMX_FLAG_ONESHOT:
		return "DVB_DMX_FLAG_ONESHOT"
	case DVB_DMX_FLAG_IMMEDIATE_START:
		return "DVB_DMX_FLAG_IMMEDIATE_START"
	default:
		return "[?? Invalid DMXFlags value]"
	}
}

func (f DMXInput) String() string {
	switch f {
	case DVB_DMX_IN_FRONTEND:
		return "DVB_DMX_IN_FRONTEND"
	case DVB_DMX_IN_DVR:
		return "DVB_DMX_IN_DVR"
	default:
		return "[?? Invalid DMXInput value]"
	}
}

func (f DMXOutput) String() string {
	switch f {
	case DVB_DMX_OUT_DECODER:
		return "DVB_DMX_OUT_DECODER"
	case DVB_DMX_OUT_TAP:
		return "DVB_DMX_OUT_TAP"
	case DVB_DMX_OUT_TS_TAP:
		return "DVB_DMX_OUT_TS_TAP"
	case DVB_DMX_OUT_TSDEMUX_TAP:
		return "DVB_DMX_OUT_TSDEMUX_TAP"
	default:
		return "[?? Invalid DMXOutput value]"
	}
}
