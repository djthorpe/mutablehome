// +build linux

/*
	Mutablehome Automation: DVB
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"path/filepath"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_PATH_WILDCARD = "/dev/dvb/adapter"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func DVBDevices() ([]uint, error) {
	if adapters, err := filepath.Glob(DVB_PATH_WILDCARD + "*"); err != nil {
		return nil, err
	} else {
		devices := make([]uint, 0, len(adapters))
		for _, file := range adapters {
			if bus, err := strconv.ParseUint(strings.TrimPrefix(file, DVB_PATH_WILDCARD), 10, 32); err == nil {
				devices = append(devices, uint(bus))
			}
		}
		return devices, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func dvb_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
