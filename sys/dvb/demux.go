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
	if file, err := os.OpenFile(DVB_FEPath(bus, demux), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func DVB_DMXStart(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_START, 0); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXStop(fd uintptr) error {
	if err := dvb_ioctl(fd, DVB_DMX_STOP, 0); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXSetFilter(fd uintptr, filter DVBDMXFilterParams) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_FILTER, 0); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_DMXSetPesFilter(fd uintptr, filter DVBDMXPesFilterParams) error {
	if err := dvb_ioctl(fd, DVB_DMX_SET_FILTER, 0); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}
