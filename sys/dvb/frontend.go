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
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
	#include <sys/ioctl.h>
	#include <linux/dvb/frontend.h>
	static int _FE_GET_INFO() { return FE_GET_INFO; }
	static int _FE_DISEQC_RESET_OVERLOAD() { return FE_DISEQC_RESET_OVERLOAD; }
	static int _FE_DISEQC_SEND_MASTER_CMD() { return FE_DISEQC_SEND_MASTER_CMD; }
	static int _FE_DISEQC_RECV_SLAVE_REPLY() { return FE_DISEQC_RECV_SLAVE_REPLY; }
	static int _FE_DISEQC_SEND_BURST() { return FE_DISEQC_SEND_BURST; }
	static int _FE_SET_TONE() { return FE_SET_TONE; }
	static int _FE_SET_VOLTAGE() { return FE_SET_VOLTAGE; }
	static int _FE_ENABLE_HIGH_LNB_VOLTAGE() { return FE_ENABLE_HIGH_LNB_VOLTAGE; }
	static int _FE_READ_STATUS() { return FE_READ_STATUS; }
	static int _FE_READ_BER() { return FE_READ_BER; }
	static int _FE_READ_SIGNAL_STRENGTH() { return FE_READ_SIGNAL_STRENGTH; }
	static int _FE_READ_SNR() { return FE_READ_SNR; }
	static int _FE_READ_UNCORRECTED_BLOCKS() { return FE_READ_UNCORRECTED_BLOCKS; }
	static int _FE_SET_FRONTEND_TUNE_MODE() { return FE_SET_FRONTEND_TUNE_MODE; }
	static int _FE_GET_EVENT() { return FE_GET_EVENT; }
	static int _FE_DISHNETWORK_SEND_LEGACY_CMD() { return FE_DISHNETWORK_SEND_LEGACY_CMD; }
	static int _FE_SET_PROPERTY() { return FE_SET_PROPERTY; }
	static int _FE_GET_PROPERTY() { return FE_GET_PROPERTY; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DVBFrontendInfo     C.struct_dvb_frontend_info
	DVBFrontendCaps     uint64
	DVBFrontendStatus   C.int
	DVBFrontendKey      uint32
	DVBFrontendValue    C.struct_dtv_property
	DVBFEDeliverySystem C.int
)

type (
	DVBFEPropertyUint32 struct {
		Key      uint32
		reserved [3]uint32
		Data     uint32
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_PATH_WILDCARD = "/dev/dvb/adapter"
)

const (
	DVB_FE_CAN_INVERSION_AUTO DVBFrontendCaps = (1 << iota)
	DVB_FE_CAN_FEC_1_2
	DVB_FE_CAN_FEC_2_3
	DVB_FE_CAN_FEC_3_4
	DVB_FE_CAN_FEC_4_5
	DVB_FE_CAN_FEC_5_6
	DVB_FE_CAN_FEC_6_7
	DVB_FE_CAN_FEC_7_8
	DVB_FE_CAN_FEC_8_9
	DVB_FE_CAN_FEC_AUTO
	DVB_FE_CAN_QPSK
	DVB_FE_CAN_QAM_16
	DVB_FE_CAN_QAM_32
	DVB_FE_CAN_QAM_64
	DVB_FE_CAN_QAM_128
	DVB_FE_CAN_QAM_256
	DVB_FE_CAN_QAM_AUTO
	DVB_FE_CAN_TRANSMISSION_MODE_AUTO
	DVB_FE_CAN_BANDWIDTH_AUTO
	DVB_FE_CAN_GUARD_INTERVAL_AUTO
	DVB_FE_CAN_HIERARCHY_AUTO
	DVB_FE_CAN_8VSB
	DVB_FE_CAN_16VSB
	DVB_FE_HAS_EXTENDED_CAPS
	DVB_FE_CAN_MULTISTREAM   DVBFrontendCaps = 0x4000000
	DVB_FE_CAN_TURBO_FEC     DVBFrontendCaps = 0x8000000
	DVB_FE_CAN_2G_MODULATION DVBFrontendCaps = 0x10000000
	DVB_FE_NEEDS_BENDING     DVBFrontendCaps = 0x20000000
	DVB_FE_CAN_RECOVER       DVBFrontendCaps = 0x40000000
	DVB_FE_CAN_MUTE_TS       DVBFrontendCaps = 0x80000000
	DVB_FE_MIN                               = DVB_FE_CAN_INVERSION_AUTO
	DVB_FE_MAX                               = DVB_FE_CAN_MUTE_TS
	DVB_FE_NONE              DVBFrontendCaps = 0
)

const (
	DVB_FE_STATUS_NONE        DVBFrontendStatus = 0x00
	DVB_FE_STATUS_HAS_SIGNAL  DVBFrontendStatus = 0x01
	DVB_FE_STATUS_HAS_CARRIER DVBFrontendStatus = 0x02
	DVB_FE_STATUS_HAS_VITERBI DVBFrontendStatus = 0x04
	DVB_FE_STATUS_HAS_SYNC    DVBFrontendStatus = 0x08
	DVB_FE_STATUS_HAS_LOCK    DVBFrontendStatus = 0x10
	DVB_FE_STATUS_TIMEDOUT    DVBFrontendStatus = 0x20
	DVB_FE_STATUS_REINIT      DVBFrontendStatus = 0x40
	DVB_FE_STATUS_MIN                           = DVB_FE_STATUS_HAS_SIGNAL
	DVB_FE_STATUS_MAX                           = DVB_FE_STATUS_REINIT
)

const (
	/* DVBv5 property Commands */
	DVB_FE_KEY_NONE               DVBFrontendKey = 0
	DVB_FE_KEY_TUNE               DVBFrontendKey = 1
	DVB_FE_KEY_CLEAR              DVBFrontendKey = 2
	DVB_FE_KEY_FREQUENCY          DVBFrontendKey = 3
	DVB_FE_KEY_MODULATION         DVBFrontendKey = 4
	DVB_FE_KEY_BANDWIDTH_HZ       DVBFrontendKey = 5
	DVB_FE_KEY_INVERSION          DVBFrontendKey = 6
	DVB_FE_KEY_DISEQC_MASTER      DVBFrontendKey = 7
	DVB_FE_KEY_SYMBOL_RATE        DVBFrontendKey = 8
	DVB_FE_KEY_INNER_FEC          DVBFrontendKey = 9
	DVB_FE_KEY_VOLTAGE            DVBFrontendKey = 10
	DVB_FE_KEY_TONE               DVBFrontendKey = 11
	DVB_FE_KEY_PILOT              DVBFrontendKey = 12
	DVB_FE_KEY_ROLLOFF            DVBFrontendKey = 13
	DVB_FE_KEY_DISEQC_SLAVE_REPLY DVBFrontendKey = 14

	/* Basic enumeration set for querying unlimited capabilities */
	DVB_FE_KEY_FE_CAPABILITY_COUNT DVBFrontendKey = 15
	DVB_FE_KEY_FE_CAPABILITY       DVBFrontendKey = 16
	DVB_FE_KEY_DELIVERY_SYSTEM     DVBFrontendKey = 17

	/* ISDB-T and ISDB-Tsb */
	DVB_FE_KEY_ISDBT_PARTIAL_RECEPTION  DVBFrontendKey = 18
	DVB_FE_KEY_ISDBT_SOUND_BROADCASTING DVBFrontendKey = 19

	DVB_FE_KEY_ISDBT_SB_SUBCHANNEL_ID DVBFrontendKey = 20
	DVB_FE_KEY_ISDBT_SB_SEGMENT_IDX   DVBFrontendKey = 21
	DVB_FE_KEY_ISDBT_SB_SEGMENT_COUNT DVBFrontendKey = 22

	DVB_FE_KEY_ISDBT_LAYERA_FEC               DVBFrontendKey = 23
	DVB_FE_KEY_ISDBT_LAYERA_MODULATION        DVBFrontendKey = 24
	DVB_FE_KEY_ISDBT_LAYERA_SEGMENT_COUNT     DVBFrontendKey = 25
	DVB_FE_KEY_ISDBT_LAYERA_TIME_INTERLEAVING DVBFrontendKey = 26

	DVB_FE_KEY_ISDBT_LAYERB_FEC               DVBFrontendKey = 27
	DVB_FE_KEY_ISDBT_LAYERB_MODULATION        DVBFrontendKey = 28
	DVB_FE_KEY_ISDBT_LAYERB_SEGMENT_COUNT     DVBFrontendKey = 29
	DVB_FE_KEY_ISDBT_LAYERB_TIME_INTERLEAVING DVBFrontendKey = 30

	DVB_FE_KEY_ISDBT_LAYERC_FEC               DVBFrontendKey = 31
	DVB_FE_KEY_ISDBT_LAYERC_MODULATION        DVBFrontendKey = 32
	DVB_FE_KEY_ISDBT_LAYERC_SEGMENT_COUNT     DVBFrontendKey = 33
	DVB_FE_KEY_ISDBT_LAYERC_TIME_INTERLEAVING DVBFrontendKey = 34

	DVB_FE_KEY_API_VERSION DVBFrontendKey = 35

	DVB_FE_KEY_CODE_RATE_HP      DVBFrontendKey = 36
	DVB_FE_KEY_CODE_RATE_LP      DVBFrontendKey = 37
	DVB_FE_KEY_GUARD_INTERVAL    DVBFrontendKey = 38
	DVB_FE_KEY_TRANSMISSION_MODE DVBFrontendKey = 39
	DVB_FE_KEY_HIERARCHY         DVBFrontendKey = 40

	DVB_FE_KEY_ISDBT_LAYER_ENABLED DVBFrontendKey = 41

	DVB_FE_KEY_STREAM_ID           DVBFrontendKey = 42
	DVB_FE_KEY_DVBT2_PLP_ID_LEGACY DVBFrontendKey = 43

	DVB_FE_KEY_ENUM_DELSYS DVBFrontendKey = 44

	/* ATSC-MH */
	DVB_FE_KEY_ATSCMH_FIC_VER           DVBFrontendKey = 45
	DVB_FE_KEY_ATSCMH_PARADE_ID         DVBFrontendKey = 46
	DVB_FE_KEY_ATSCMH_NOG               DVBFrontendKey = 47
	DVB_FE_KEY_ATSCMH_TNOG              DVBFrontendKey = 48
	DVB_FE_KEY_ATSCMH_SGN               DVBFrontendKey = 49
	DVB_FE_KEY_ATSCMH_PRC               DVBFrontendKey = 50
	DVB_FE_KEY_ATSCMH_RS_FRAME_MODE     DVBFrontendKey = 51
	DVB_FE_KEY_ATSCMH_RS_FRAME_ENSEMBLE DVBFrontendKey = 52
	DVB_FE_KEY_ATSCMH_RS_CODE_MODE_PRI  DVBFrontendKey = 53
	DVB_FE_KEY_ATSCMH_RS_CODE_MODE_SEC  DVBFrontendKey = 54
	DVB_FE_KEY_ATSCMH_SCCC_BLOCK_MODE   DVBFrontendKey = 55
	DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_A  DVBFrontendKey = 56
	DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_B  DVBFrontendKey = 57
	DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_C  DVBFrontendKey = 58
	DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_D  DVBFrontendKey = 59

	DVB_FE_KEY_INTERLEAVING DVBFrontendKey = 60
	DVB_FE_KEY_LNA          DVBFrontendKey = 61

	/* Quality parameters */
	DVB_FE_KEY_STAT_SIGNAL_STRENGTH      DVBFrontendKey = 62
	DVB_FE_KEY_STAT_CNR                  DVBFrontendKey = 63
	DVB_FE_KEY_STAT_PRE_ERROR_BIT_COUNT  DVBFrontendKey = 64
	DVB_FE_KEY_STAT_PRE_TOTAL_BIT_COUNT  DVBFrontendKey = 65
	DVB_FE_KEY_STAT_POST_ERROR_BIT_COUNT DVBFrontendKey = 66
	DVB_FE_KEY_STAT_POST_TOTAL_BIT_COUNT DVBFrontendKey = 67
	DVB_FE_KEY_STAT_ERROR_BLOCK_COUNT    DVBFrontendKey = 68
	DVB_FE_KEY_STAT_TOTAL_BLOCK_COUNT    DVBFrontendKey = 69

	/* Physical layer scrambling */
	DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX DVBFrontendKey = 70

	// Minimum and maximum
	DVB_FE_KEY_MIN = DVB_FE_KEY_TUNE
	DVB_FE_KEY_MAX = DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX
)

const (
	DVB_FE_SYS_NONE         DVBFEDeliverySystem = 0
	DVB_FE_SYS_DVBC_ANNEX_A DVBFEDeliverySystem = iota // Cable TV: DVB-C following ITU-T J.83 Annex A spec
	DVB_FE_SYS_DVBC_ANNEX_B                            // Cable TV: DVB-C following ITU-T J.83 Annex B spec (ClearQAM)
	DVB_FE_SYS_DVBT                                    // Terrestrial TV: DVB-T
	DVB_FE_SYS_DSS                                     // Satellite TV: DSS (not fully supported)
	DVB_FE_SYS_DVBS                                    // Satellite TV: DVB-S
	DVB_FE_SYS_DVBS2                                   // Satellite TV: DVB-S2
	DVB_FE_SYS_DVBH                                    // Terrestrial TV (mobile): DVB-H (standard deprecated)
	DVB_FE_SYS_ISDBT                                   // Terrestrial TV: ISDB-T
	DVB_FE_SYS_ISDBS                                   // Satellite TV: ISDB-S
	DVB_FE_SYS_ISDBC                                   // Cable TV: ISDB-C (no drivers yet)
	DVB_FE_SYS_ATSC                                    // Terrestrial TV: ATSC
	DVB_FE_SYS_ATSCMH                                  // Terrestrial TV (mobile): ATSC-M/H
	DVB_FE_SYS_DTMB                                    // Terrestrial TV: DTMB
	DVB_FE_SYS_CMMB                                    // Terrestrial TV (mobile): CMMB (not fully supported)
	DVB_FE_SYS_DAB                                     // Digital audio: DAB (not fully supported)
	DVB_FE_SYS_DVBT2                                   // Terrestrial TV: DVB-T2
	DVB_FE_SYS_TURBO                                   // Satellite TV: DVB-S Turbo
	DVB_FE_SYS_DVBC_ANNEX_C                            // Cable TV: DVB-C following ITU-T J.83 Annex C spec
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	DVB_FE_GET_INFO     = uintptr(C._FE_GET_INFO())
	DVB_FE_READ_STATUS  = uintptr(C._FE_READ_STATUS())
	DVB_FE_GET_PROPERTY = uintptr(C._FE_GET_PROPERTY())
	DVB_FE_SET_PROPERTY = uintptr(C._FE_SET_PROPERTY())
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
// PUBLIC METHODS: FRONT END

func DVB_FEPath(bus uint) string {
	return fmt.Sprintf("%v%v/frontend%v", DVB_PATH_WILDCARD, bus, bus)
}

func DVB_FEOpen(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_FEPath(bus), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func DVB_FEGetInfo(fd uintptr) (DVBFrontendInfo, error) {
	var info DVBFrontendInfo
	if err := dvb_ioctl(fd, DVB_FE_GET_INFO, unsafe.Pointer(&info)); err != 0 {
		return info, os.NewSyscallError("dvb_ioctl", err)
	} else {
		return info, nil
	}
}

func DVB_FEReadStatus(fd uintptr) (DVBFrontendStatus, error) {
	var status DVBFrontendStatus
	if err := dvb_ioctl(fd, DVB_FE_READ_STATUS, unsafe.Pointer(&status)); err != 0 {
		return status, os.NewSyscallError("dvb_ioctl", err)
	} else {
		return status, nil
	}
}

func DVB_FEGetPropertyUint32(fd uintptr, key DVBFrontendKey) (uint32, error) {
	property := DVBFEPropertyUint32{Key: uint32(key)}
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, DVB_FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return 0, os.NewSyscallError("dvb_ioctl", err)
	} else {
		return property.Data, nil
	}
}

func DVB_FESetPropertyUint32(fd uintptr, key DVBFrontendKey, value uint32) error {
	property := DVBFEPropertyUint32{Key: uint32(key), Data: value}
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, DVB_FE_SET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func DVB_FEVersion(fd uintptr) (uint, uint, error) {
	if version, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_API_VERSION); err != nil {
		return 0, 0, err
	} else {
		major := version >> 8 & 0xFF
		minor := version & 0xFF
		return uint(major), uint(minor), nil
	}
}

func DVB_FEDeliverySystem(fd uintptr) (DVBFEDeliverySystem, error) {
	if sys, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_DELIVERY_SYSTEM); err != nil {
		return DVB_FE_SYS_NONE, err
	} else {
		return DVBFEDeliverySystem(sys), nil
	}
}

func DVB_FESetDeliverySystem(fd uintptr, sys DVBFEDeliverySystem) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_DELIVERY_SYSTEM, uint32(sys))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this DVBFrontendInfo) Name() string {
	return C.GoString(&this.name[0])
}

func (this DVBFrontendInfo) FrequencyMin() uint32 {
	return uint32(this.frequency_min)
}

func (this DVBFrontendInfo) FrequencyMax() uint32 {
	return uint32(this.frequency_max)
}

func (this DVBFrontendInfo) FrequencyStepSize() uint32 {
	return uint32(this.frequency_stepsize)
}

func (this DVBFrontendInfo) FrequencyTolerance() uint32 {
	return uint32(this.frequency_tolerance)
}

func (this DVBFrontendInfo) SymbolrateMin() uint32 {
	return uint32(this.symbol_rate_min)
}

func (this DVBFrontendInfo) SymbolrateMax() uint32 {
	return uint32(this.symbol_rate_max)
}

func (this DVBFrontendInfo) SymbolrateTolerance() uint32 {
	return uint32(this.symbol_rate_max)
}

func (this DVBFrontendInfo) Caps() DVBFrontendCaps {
	return DVBFrontendCaps(this.caps)
}

func (this DVBFrontendInfo) String() string {
	return "<DVBFrontendInfo" +
		" name=" + strconv.Quote(this.Name()) +
		" caps=" + fmt.Sprint(this.Caps()) +
		" frequency=" + fmt.Sprintf("{ %v,%v }", this.FrequencyMin(), this.FrequencyMax()) +
		" symbolrate=" + fmt.Sprintf("{ %v,%v }", this.SymbolrateMin(), this.SymbolrateMax()) +
		">"
}

func (f DVBFrontendCaps) String() string {
	str := ""
	if f == DVB_FE_NONE {
		return f.StringFlag()
	}
	for v := DVB_FE_MIN; v <= DVB_FE_MAX; v = v << 1 {
		if v&f == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (v DVBFrontendCaps) StringFlag() string {
	switch v {
	case DVB_FE_NONE:
		return "DVB_FE_NONE"
	case DVB_FE_CAN_INVERSION_AUTO:
		return "DVB_FE_CAN_INVERSION_AUTO"
	case DVB_FE_CAN_FEC_1_2:
		return "DVB_FE_CAN_FEC_1_2"
	case DVB_FE_CAN_FEC_2_3:
		return "DVB_FE_CAN_FEC_2_3"
	case DVB_FE_CAN_FEC_3_4:
		return "DVB_FE_CAN_FEC_3_4"
	case DVB_FE_CAN_FEC_4_5:
		return "DVB_FE_CAN_FEC_4_5"
	case DVB_FE_CAN_FEC_5_6:
		return "DVB_FE_CAN_FEC_5_6"
	case DVB_FE_CAN_FEC_6_7:
		return "DVB_FE_CAN_FEC_6_7"
	case DVB_FE_CAN_FEC_7_8:
		return "DVB_FE_CAN_FEC_7_8"
	case DVB_FE_CAN_FEC_8_9:
		return "DVB_FE_CAN_FEC_8_9"
	case DVB_FE_CAN_FEC_AUTO:
		return "DVB_FE_CAN_FEC_AUTO"
	case DVB_FE_CAN_QPSK:
		return "DVB_FE_CAN_QPSK"
	case DVB_FE_CAN_QAM_16:
		return "DVB_FE_CAN_QAM_16"
	case DVB_FE_CAN_QAM_32:
		return "DVB_FE_CAN_QAM_32"
	case DVB_FE_CAN_QAM_64:
		return "DVB_FE_CAN_QAM_64"
	case DVB_FE_CAN_QAM_128:
		return "DVB_FE_CAN_QAM_128"
	case DVB_FE_CAN_QAM_256:
		return "DVB_FE_CAN_QAM_256"
	case DVB_FE_CAN_QAM_AUTO:
		return "DVB_FE_CAN_QAM_AUTO"
	case DVB_FE_CAN_TRANSMISSION_MODE_AUTO:
		return "DVB_FE_CAN_TRANSMISSION_MODE_AUTO"
	case DVB_FE_CAN_BANDWIDTH_AUTO:
		return "DVB_FE_CAN_BANDWIDTH_AUTO"
	case DVB_FE_CAN_GUARD_INTERVAL_AUTO:
		return "DVB_FE_CAN_GUARD_INTERVAL_AUTO"
	case DVB_FE_CAN_HIERARCHY_AUTO:
		return "DVB_FE_CAN_HIERARCHY_AUTO"
	case DVB_FE_CAN_8VSB:
		return "DVB_FE_CAN_8VSB"
	case DVB_FE_CAN_16VSB:
		return "DVB_FE_CAN_16VSB"
	case DVB_FE_HAS_EXTENDED_CAPS:
		return "DVB_FE_HAS_EXTENDED_CAPS"
	case DVB_FE_CAN_MULTISTREAM:
		return "DVB_FE_CAN_MULTISTREAM"
	case DVB_FE_CAN_TURBO_FEC:
		return "DVB_FE_CAN_TURBO_FEC"
	case DVB_FE_CAN_2G_MODULATION:
		return "DVB_FE_CAN_2G_MODULATION"
	case DVB_FE_NEEDS_BENDING:
		return "DVB_FE_NEEDS_BENDING"
	case DVB_FE_CAN_RECOVER:
		return "DVB_FE_CAN_RECOVER"
	case DVB_FE_CAN_MUTE_TS:
		return "DVB_FE_CAN_MUTE_TS"
	default:
		return "[?? Invalid DVBFrontendCaps value]"
	}
}

func (f DVBFrontendStatus) String() string {
	str := ""
	if f == DVB_FE_STATUS_NONE {
		return f.StringFlag()
	}
	for v := DVB_FE_STATUS_MIN; v <= DVB_FE_STATUS_MAX; v = v << 1 {
		if v&f == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (s DVBFrontendStatus) StringFlag() string {
	switch s {
	case DVB_FE_STATUS_NONE:
		return "DVB_FE_STATUS_NONE"
	case DVB_FE_STATUS_HAS_SIGNAL:
		return "DVB_FE_STATUS_HAS_SIGNAL"
	case DVB_FE_STATUS_HAS_CARRIER:
		return "DVB_FE_STATUS_HAS_CARRIER"
	case DVB_FE_STATUS_HAS_VITERBI:
		return "DVB_FE_STATUS_HAS_VITERBI"
	case DVB_FE_STATUS_HAS_SYNC:
		return "DVB_FE_STATUS_HAS_SYNC"
	case DVB_FE_STATUS_HAS_LOCK:
		return "DVB_FE_STATUS_HAS_LOCK"
	case DVB_FE_STATUS_TIMEDOUT:
		return "DVB_FE_STATUS_TIMEDOUT"
	case DVB_FE_STATUS_REINIT:
		return "DVB_FE_STATUS_REINIT"
	default:
		return "[?? Invalid DVBFrontendStatus value]"
	}
}

func (v DVBFEDeliverySystem) String() string {
	switch v {
	case DVB_FE_SYS_NONE:
		return "DVB_FE_SYS_NONE"
	case DVB_FE_SYS_DVBC_ANNEX_A:
		return "DVB_FE_SYS_DVBC_ANNEX_A"
	case DVB_FE_SYS_DVBC_ANNEX_B:
		return "DVB_FE_SYS_DVBC_ANNEX_B"
	case DVB_FE_SYS_DVBT:
		return "DVB_FE_SYS_DVBT"
	case DVB_FE_SYS_DSS:
		return "DVB_FE_SYS_DSS"
	case DVB_FE_SYS_DVBS:
		return "DVB_FE_SYS_DVBS"
	case DVB_FE_SYS_DVBS2:
		return "DVB_FE_SYS_DVBS2"
	case DVB_FE_SYS_DVBH:
		return "DVB_FE_SYS_DVBH"
	case DVB_FE_SYS_ISDBT:
		return "DVB_FE_SYS_ISDBT"
	case DVB_FE_SYS_ISDBS:
		return "DVB_FE_SYS_ISDBS"
	case DVB_FE_SYS_ISDBC:
		return "DVB_FE_SYS_ISDBC"
	case DVB_FE_SYS_ATSC:
		return "DVB_FE_SYS_ATSC"
	case DVB_FE_SYS_ATSCMH:
		return "DVB_FE_SYS_ATSCMH"
	case DVB_FE_SYS_DTMB:
		return "DVB_FE_SYS_DTMB"
	case DVB_FE_SYS_CMMB:
		return "DVB_FE_SYS_CMMB"
	case DVB_FE_SYS_DAB:
		return "DVB_FE_SYS_DAB"
	case DVB_FE_SYS_DVBT2:
		return "DVB_FE_SYS_DVBT2"
	case DVB_FE_SYS_TURBO:
		return "DVB_FE_SYS_TURBO"
	case DVB_FE_SYS_DVBC_ANNEX_C:
		return "DVB_FE_SYS_DVBC_ANNEX_C"
	default:
		return "[?? Invalid DVBFEDeliverySystem value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func dvb_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
