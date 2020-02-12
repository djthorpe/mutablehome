// +build linux

/*
	Mutablehome Automation: DVB
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"fmt"
	"os"
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
		Timeout uint32
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
	DMX_PES_AUDIO DMXStreamType = iota
	DVB_DMX_PES_VIDEO
	DVB_DMX_PES_TELETEXT
	DVB_DMX_PES_SUBTITLE
	DVB_DMX_PES_PCR
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
	DVB_DMX_FLAG_NONE DMXFlags = iota
	DVB_DMX_FLAG_CHECK_CRC
	DVB_DMX_FLAG_ONESHOT
	DVB_DMX_FLAG_IMMEDIATE_START
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

func DVB_DMXOpen(bus, demux uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_DMXPath(bus, demux), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func DVB_DMXStart(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_START, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXStop(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_STOP, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXSetBufferSize(fd uintptr, size uint32) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_BUFFER_SIZE, unsafe.Pointer(uintptr(size))); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXSetSectionFilter(fd uintptr, filter DMXSectionFilter) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_FILTER, unsafe.Pointer(&filter)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXSetStreamFilter(fd uintptr, filter DMXStreamFilter) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_PES_FILTER, unsafe.Pointer(&filter)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXAddPid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DVB_DMX_ADD_PID, unsafe.Pointer(uintptr(pid))); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXRemovePid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DVB_DMX_REMOVE_PID, unsafe.Pointer(uintptr(pid))); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXGetStreamPids(fd uintptr) (map[DMXStreamType]uint16, error) {
	var pids [5]uint16
	if err := dvb_ioctl(fd, DVB_DMX_GET_PES_PIDS, unsafe.Pointer(&pids)); err != 0 {
		return nil, os.NewSyscallError("dvb_ioctl", err)
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

func (s DMXStreamType) String() string {
	switch s {
	case DMX_PES_AUDIO:
		return "DMX_PES_AUDIO"
	case DVB_DMX_PES_VIDEO:
		return "DVB_DMX_PES_VIDEO"
	case DVB_DMX_PES_TELETEXT:
		return "DVB_DMX_PES_TELETEXT"
	case DVB_DMX_PES_SUBTITLE:
		return "DVB_DMX_PES_SUBTITLE"
	case DVB_DMX_PES_PCR:
		return "DVB_DMX_PES_PCR"
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
		return "[?? Invalid DMXStreamType]"
	}
}
