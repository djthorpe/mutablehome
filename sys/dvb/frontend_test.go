// +build linux

/*
	Mutablehome Automation: DVB
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package dvb_test

import (
	"testing"

	// Frameworks
	dvb "github.com/djthorpe/mutablehome/sys/dvb"
)

func Test_Frontend_000(t *testing.T) {
	t.Log("Test_Frontend_000")
}

func Test_Frontend_001(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		t.Log("devices=", devices)
	}
}

func Test_Frontend_002(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_FEOpen(device); err != nil {
				t.Error(err)
			} else if info, err := dvb.DVB_FEGetInfo(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", info)
				dev.Close()
			}
		}
	}
}

func Test_Frontend_003(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_FEOpen(device); err != nil {
				t.Error(err)
			} else if status, err := dvb.DVB_FEReadStatus(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", status)
				dev.Close()
			}
		}
	}
}

func Test_Frontend_004(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_FEOpen(device); err != nil {
				t.Error(err)
			} else if major, minor, err := dvb.DVB_FEVersion(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", major, minor)
				dev.Close()
			}
		}
	}
}

func Test_Frontend_005(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_FEOpen(device); err != nil {
				t.Error(err)
			} else if sys, err := dvb.DVB_FEDeliverySystem(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", sys)
				dev.Close()
			}
		}
	}
}
