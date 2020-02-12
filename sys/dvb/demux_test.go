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

func Test_Demux_000(t *testing.T) {
	t.Log("Test_Demux_000")
}

func Test_Demux_001(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		t.Log("devices=", devices)
	}
}

func Test_Demux_002(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_DMXOpen(device, 0); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", dev.Fd())
				dev.Close()
			}
		}
	}
}

func Test_Demux_003(t *testing.T) {
	if devices, err := dvb.DVBDevices(); err != nil {
		t.Error(err)
	} else {
		for _, device := range devices {
			if dev, err := dvb.DVB_DMXOpen(device, 0); err != nil {
				t.Error(err)
			} else if err := dvb.DVB_DMXSetStreamFilter(dev.Fd(), dvb.DMXStreamFilter{}); err != nil {
				t.Error(err)
			} else if err := dvb.DVB_DMXStart(dev.Fd()); err != nil {
				t.Error(err)
			} else if pids, err := dvb.DVB_DMXGetStreamPids(dev.Fd()); err != nil {
				t.Error(err)
			} else if err := dvb.DVB_DMXStop(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log(dev.Name(), "=>", pids)
				dev.Close()
			}
		}
	}
}
