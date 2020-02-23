// +build linux

/*
	Mutablehome Automation: DVB
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	// Frameworks

	mutablehome "github.com/djthorpe/mutablehome"
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
	DVBFrontendInfo   C.struct_dvb_frontend_info
	DVBFrontendCaps   uint64
	DVBFrontendStatus C.int
	DVBFrontendKey    uint32
	DVBFrontendValue  C.struct_dtv_property
	DVBFrontendScale  uint8
)

type (
	DVBFEPropertyUint32 struct {
		Key      DVBFrontendKey
		reserved [3]uint32
		Data     uint32
		result   C.int
	}
	DVBFEPropertyEnum struct {
		Key       DVBFrontendKey
		reserved  [3]uint32
		Data      [32]uint8
		Len       uint32
		reserved1 [11]byte
		reserved2 uintptr
		result    C.int
	}
	DVBFrontendStats struct {
		Key       DVBFrontendKey
		reserved  [3]uint32
		Len       uint8
		Data      [9 * 4]byte // Up to four stats of 9 bytes
		reserved1 [11]byte
		reserved2 uintptr
		result    C.int
	}
	DVBFrontendStat struct {
		Scale DVBFrontendScale
		Value int64
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

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
	DVB_FE_STAT_SIGNAL_STRENGTH      DVBFrontendKey = 62
	DVB_FE_STAT_CNR                  DVBFrontendKey = 63
	DVB_FE_STAT_PRE_ERROR_BIT_COUNT  DVBFrontendKey = 64
	DVB_FE_STAT_PRE_TOTAL_BIT_COUNT  DVBFrontendKey = 65
	DVB_FE_STAT_POST_ERROR_BIT_COUNT DVBFrontendKey = 66
	DVB_FE_STAT_POST_TOTAL_BIT_COUNT DVBFrontendKey = 67
	DVB_FE_STAT_ERROR_BLOCK_COUNT    DVBFrontendKey = 68
	DVB_FE_STAT_TOTAL_BLOCK_COUNT    DVBFrontendKey = 69

	/* Physical layer scrambling */
	DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX DVBFrontendKey = 70

	// Minimum and maximum
	DVB_FE_KEY_MIN = DVB_FE_KEY_TUNE
	DVB_FE_KEY_MAX = DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX
)

const (
	DVB_FE_SCALE_NONE DVBFrontendScale = iota
	DVB_FE_SCALE_DECIBEL
	DVB_FE_SCALE_RELATIVE
	DVB_FE_SCALE_COUNTER
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
// PUBLIC METHODS: FRONT END

func DVB_FEPath(bus, frontend uint) string {
	return fmt.Sprintf("%v%v/frontend%v", DVB_PATH_WILDCARD, bus, frontend)
}

func DVB_FEOpen(bus, frontend uint) (*os.File, error) {
	if file, err := os.OpenFile(DVB_FEPath(bus, frontend), os.O_SYNC|os.O_RDWR, 0); err != nil {
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
	property := DVBFEPropertyUint32{Key: key}
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
	property := DVBFEPropertyUint32{Key: key, Data: value}
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, DVB_FE_SET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return os.NewSyscallError("dvb_ioctl", err)
	} else {
		return nil
	}
}

func DVB_FEGetPropertyEnum(fd uintptr, key DVBFrontendKey) ([]uint8, error) {
	property := DVBFEPropertyEnum{Key: key}
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, DVB_FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return nil, os.NewSyscallError("dvb_ioctl", err)
	} else {
		return property.Data[0:property.Len], nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: FRONT END GET/SET PROPERTIES

func DVB_FEVersion(fd uintptr) (uint, uint, error) {
	if version, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_API_VERSION); err != nil {
		return 0, 0, err
	} else {
		major := version >> 8 & 0xFF
		minor := version & 0xFF
		return uint(major), uint(minor), nil
	}
}

func DVB_FETune(fd uintptr) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_TUNE, uint32(0))
}

func DVB_FEClear(fd uintptr) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_CLEAR, uint32(0))
}

func DVB_FEDeliverySystemEnum(fd uintptr) ([]mutablehome.DVBDeliverySystem, error) {
	if data, err := DVB_FEGetPropertyEnum(fd, DVB_FE_KEY_ENUM_DELSYS); err != nil {
		return nil, err
	} else {
		enum := make([]mutablehome.DVBDeliverySystem, len(data))
		for i, value := range data {
			enum[i] = mutablehome.DVBDeliverySystem(value)
		}
		return enum, nil
	}
}

func DVB_FEDeliverySystem(fd uintptr) (mutablehome.DVBDeliverySystem, error) {
	if sys, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_DELIVERY_SYSTEM); err != nil {
		return mutablehome.DVB_SYS_NONE, err
	} else {
		return mutablehome.DVBDeliverySystem(sys), nil
	}
}

func DVB_FESetDeliverySystem(fd uintptr, sys mutablehome.DVBDeliverySystem) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_DELIVERY_SYSTEM, uint32(sys))
}

func DVB_FEFrequency(fd uintptr) (uint, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_FREQUENCY); err != nil {
		return 0, err
	} else {
		return uint(value), err
	}
}

func DVB_FESetFrequency(fd uintptr, value uint) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_FREQUENCY, uint32(value))
}

func DVB_FEBandwidth(fd uintptr) (uint, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_BANDWIDTH_HZ); err != nil {
		return 0, err
	} else {
		return uint(value), err
	}
}

func DVB_FESetBandwidth(fd uintptr, value uint) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_BANDWIDTH_HZ, uint32(value))
}

func DVB_FEModulation(fd uintptr) (mutablehome.DVBModulation, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_MODULATION); err != nil {
		return 0, err
	} else {
		return mutablehome.DVBModulation(value), err
	}
}

func DVB_FESetModulation(fd uintptr, value mutablehome.DVBModulation) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_MODULATION, uint32(value))
}

func DVB_FEInversion(fd uintptr) (mutablehome.DVBInversion, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_INVERSION); err != nil {
		return 0, err
	} else {
		return mutablehome.DVBInversion(value), err
	}
}

func DVB_FESetInversion(fd uintptr, value mutablehome.DVBInversion) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_INVERSION, uint32(value))
}

func DVB_FEGuardInterval(fd uintptr) (mutablehome.DVBGuardInterval, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_GUARD_INTERVAL); err != nil {
		return 0, err
	} else {
		return mutablehome.DVBGuardInterval(value), err
	}
}

func DVB_FESetGuardInterval(fd uintptr, value mutablehome.DVBGuardInterval) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_GUARD_INTERVAL, uint32(value))
}

func DVB_FEHierarchy(fd uintptr) (mutablehome.DVBHierarchy, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_HIERARCHY); err != nil {
		return mutablehome.DVB_HIERARCHY_NONE, err
	} else {
		return mutablehome.DVBHierarchy(value), err
	}
}

func DVB_FESetHierarchy(fd uintptr, value mutablehome.DVBHierarchy) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_HIERARCHY, uint32(value))
}

func DVB_FEInnerFEC(fd uintptr) (mutablehome.DVBCodeRate, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_INNER_FEC); err != nil {
		return mutablehome.DVB_FEC_NONE, err
	} else {
		return mutablehome.DVBCodeRate(value), err
	}
}

func DVB_FESetInnerFEC(fd uintptr, value mutablehome.DVBCodeRate) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_INNER_FEC, uint32(value))
}

func DVB_FECodeRateLP(fd uintptr) (mutablehome.DVBCodeRate, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_CODE_RATE_LP); err != nil {
		return mutablehome.DVB_FEC_NONE, err
	} else {
		return mutablehome.DVBCodeRate(value), err
	}
}

func DVB_FESetCodeRateLP(fd uintptr, value mutablehome.DVBCodeRate) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_CODE_RATE_LP, uint32(value))
}

func DVB_FECodeRateHP(fd uintptr) (mutablehome.DVBCodeRate, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_CODE_RATE_HP); err != nil {
		return mutablehome.DVB_FEC_NONE, err
	} else {
		return mutablehome.DVBCodeRate(value), err
	}
}

func DVB_FESetCodeRateHP(fd uintptr, value mutablehome.DVBCodeRate) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_CODE_RATE_HP, uint32(value))
}

func DVB_FETransmitMode(fd uintptr) (mutablehome.DVBTransmitMode, error) {
	if value, err := DVB_FEGetPropertyUint32(fd, DVB_FE_KEY_TRANSMISSION_MODE); err != nil {
		return 0, err
	} else {
		return mutablehome.DVBTransmitMode(value), err
	}
}

func DVB_FESetTransmitMode(fd uintptr, value mutablehome.DVBTransmitMode) error {
	return DVB_FESetPropertyUint32(fd, DVB_FE_KEY_TRANSMISSION_MODE, uint32(value))
}

func DVB_FEStats(fd uintptr) (map[DVBFrontendKey]DVBFrontendStat, error) {
	stats := [...]DVBFrontendStats{
		DVBFrontendStats{Key: DVB_FE_STAT_SIGNAL_STRENGTH, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_CNR, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_PRE_ERROR_BIT_COUNT, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_PRE_TOTAL_BIT_COUNT, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_POST_ERROR_BIT_COUNT, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_POST_TOTAL_BIT_COUNT, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_ERROR_BLOCK_COUNT, Len: 4},
		DVBFrontendStats{Key: DVB_FE_STAT_TOTAL_BLOCK_COUNT, Len: 4},
	}
	properties := C.struct_dtv_properties{
		C.uint(len(stats)), (*C.struct_dtv_property)(unsafe.Pointer(&stats[0])),
	}
	if err := dvb_ioctl(fd, DVB_FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return nil, os.NewSyscallError("dvb_ioctl", err)
	} else {
		statMap := make(map[DVBFrontendKey]DVBFrontendStat, len(stats))
		for _, s := range stats {
			if s.Len > 0 {
				stat := DVBFrontendStat{}
				if err := binary.Read(bytes.NewReader(s.Data[:]), binary.LittleEndian, &stat); err != nil {
					return nil, err
				} else if stat.Scale != 0 {
					statMap[s.Key] = stat
				}
			}
		}
		return statMap, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////
// DVBFrontendStat

func (stat DVBFrontendStat) Decibel() float64 {
	return float64(stat.Value) * 0.001
}

func (stat DVBFrontendStat) Relative() float64 {
	return float64(stat.Value&0xFFFF) * 100
}

func (stat DVBFrontendStat) Counter() uint64 {
	return uint64(stat.Value)
}

func (stat DVBFrontendStat) String() string {
	switch stat.Scale {
	case DVB_FE_SCALE_COUNTER:
		return fmt.Sprint(stat.Counter())
	case DVB_FE_SCALE_DECIBEL:
		return fmt.Sprintf("%.1fdB", stat.Decibel())
	case DVB_FE_SCALE_RELATIVE:
		return fmt.Sprintf("%.1f%%", stat.Relative())
	default:
		return fmt.Sprintf("{scale=0x%02X,value=0x%08X}", uint8(stat.Scale), stat.Value)
	}
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

func (s DVBFrontendScale) String() string {
	switch s {
	case DVB_FE_SCALE_NONE:
		return "DVB_FE_SCALE_NONE"
	case DVB_FE_SCALE_DECIBEL:
		return "DVB_FE_SCALE_DECIBEL"
	case DVB_FE_SCALE_RELATIVE:
		return "DVB_FE_SCALE_RELATIVE"
	case DVB_FE_SCALE_COUNTER:
		return "DVB_FE_SCALE_COUNTER"
	default:
		return "[?? Invalid DVBFrontendScale value]"
	}
}

func (k DVBFrontendKey) String() string {
	switch k {
	case DVB_FE_KEY_NONE:
		return "DVB_FE_KEY_NONE"
	case DVB_FE_KEY_TUNE:
		return "DVB_FE_KEY_TUNE"
	case DVB_FE_KEY_CLEAR:
		return "DVB_FE_KEY_CLEAR"
	case DVB_FE_KEY_FREQUENCY:
		return "DVB_FE_KEY_FREQUENCY"
	case DVB_FE_KEY_MODULATION:
		return "DVB_FE_KEY_MODULATION"
	case DVB_FE_KEY_BANDWIDTH_HZ:
		return "DVB_FE_KEY_BANDWIDTH_HZ"
	case DVB_FE_KEY_INVERSION:
		return "DVB_FE_KEY_INVERSION"
	case DVB_FE_KEY_DISEQC_MASTER:
		return "DVB_FE_KEY_DISEQC_MASTER"
	case DVB_FE_KEY_SYMBOL_RATE:
		return "DVB_FE_KEY_SYMBOL_RATE"
	case DVB_FE_KEY_INNER_FEC:
		return "DVB_FE_KEY_INNER_FEC"
	case DVB_FE_KEY_VOLTAGE:
		return "DVB_FE_KEY_VOLTAGE"
	case DVB_FE_KEY_TONE:
		return "DVB_FE_KEY_TONE"
	case DVB_FE_KEY_PILOT:
		return "DVB_FE_KEY_PILOT"
	case DVB_FE_KEY_ROLLOFF:
		return "DVB_FE_KEY_ROLLOFF"
	case DVB_FE_KEY_DISEQC_SLAVE_REPLY:
		return "DVB_FE_KEY_DISEQC_SLAVE_REPLY"
	case DVB_FE_KEY_FE_CAPABILITY_COUNT:
		return "DVB_FE_KEY_FE_CAPABILITY_COUNT"
	case DVB_FE_KEY_FE_CAPABILITY:
		return "DVB_FE_KEY_FE_CAPABILITY"
	case DVB_FE_KEY_DELIVERY_SYSTEM:
		return "DVB_FE_KEY_DELIVERY_SYSTEM"
	case DVB_FE_KEY_ISDBT_PARTIAL_RECEPTION:
		return "DVB_FE_KEY_ISDBT_PARTIAL_RECEPTION"
	case DVB_FE_KEY_ISDBT_SOUND_BROADCASTING:
		return "DVB_FE_KEY_ISDBT_SOUND_BROADCASTING"
	case DVB_FE_KEY_ISDBT_SB_SUBCHANNEL_ID:
		return "DVB_FE_KEY_ISDBT_SB_SUBCHANNEL_ID"
	case DVB_FE_KEY_ISDBT_SB_SEGMENT_IDX:
		return "DVB_FE_KEY_ISDBT_SB_SEGMENT_IDX"
	case DVB_FE_KEY_ISDBT_SB_SEGMENT_COUNT:
		return "DVB_FE_KEY_ISDBT_SB_SEGMENT_COUNT"
	case DVB_FE_KEY_ISDBT_LAYERA_FEC:
		return "DVB_FE_KEY_ISDBT_LAYERA_FEC"
	case DVB_FE_KEY_ISDBT_LAYERA_MODULATION:
		return "DVB_FE_KEY_ISDBT_LAYERA_MODULATION"
	case DVB_FE_KEY_ISDBT_LAYERA_SEGMENT_COUNT:
		return "DVB_FE_KEY_ISDBT_LAYERA_SEGMENT_COUNT"
	case DVB_FE_KEY_ISDBT_LAYERA_TIME_INTERLEAVING:
		return "DVB_FE_KEY_ISDBT_LAYERA_TIME_INTERLEAVING"
	case DVB_FE_KEY_ISDBT_LAYERB_FEC:
		return "DVB_FE_KEY_ISDBT_LAYERB_FEC"
	case DVB_FE_KEY_ISDBT_LAYERB_MODULATION:
		return "DVB_FE_KEY_ISDBT_LAYERB_MODULATION"
	case DVB_FE_KEY_ISDBT_LAYERB_SEGMENT_COUNT:
		return "DVB_FE_KEY_ISDBT_LAYERB_SEGMENT_COUNT"
	case DVB_FE_KEY_ISDBT_LAYERB_TIME_INTERLEAVING:
		return "DVB_FE_KEY_ISDBT_LAYERB_TIME_INTERLEAVING"
	case DVB_FE_KEY_ISDBT_LAYERC_FEC:
		return "DVB_FE_KEY_ISDBT_LAYERC_FEC"
	case DVB_FE_KEY_ISDBT_LAYERC_MODULATION:
		return "DVB_FE_KEY_ISDBT_LAYERC_MODULATION"
	case DVB_FE_KEY_ISDBT_LAYERC_SEGMENT_COUNT:
		return "DVB_FE_KEY_ISDBT_LAYERC_SEGMENT_COUNT"
	case DVB_FE_KEY_ISDBT_LAYERC_TIME_INTERLEAVING:
		return "DVB_FE_KEY_ISDBT_LAYERC_TIME_INTERLEAVING"
	case DVB_FE_KEY_API_VERSION:
		return "DVB_FE_KEY_API_VERSION"
	case DVB_FE_KEY_CODE_RATE_HP:
		return "DVB_FE_KEY_CODE_RATE_HP"
	case DVB_FE_KEY_CODE_RATE_LP:
		return "DVB_FE_KEY_CODE_RATE_LP"
	case DVB_FE_KEY_GUARD_INTERVAL:
		return "DVB_FE_KEY_GUARD_INTERVAL"
	case DVB_FE_KEY_TRANSMISSION_MODE:
		return "DVB_FE_KEY_TRANSMISSION_MODE"
	case DVB_FE_KEY_HIERARCHY:
		return "DVB_FE_KEY_HIERARCHY"
	case DVB_FE_KEY_ISDBT_LAYER_ENABLED:
		return "DVB_FE_KEY_ISDBT_LAYER_ENABLED"
	case DVB_FE_KEY_STREAM_ID:
		return "DVB_FE_KEY_STREAM_ID"
	case DVB_FE_KEY_DVBT2_PLP_ID_LEGACY:
		return "DVB_FE_KEY_DVBT2_PLP_ID_LEGACY"
	case DVB_FE_KEY_ENUM_DELSYS:
		return "DVB_FE_KEY_ENUM_DELSYS"
	case DVB_FE_KEY_ATSCMH_FIC_VER:
		return "DVB_FE_KEY_ATSCMH_FIC_VER"
	case DVB_FE_KEY_ATSCMH_PARADE_ID:
		return "DVB_FE_KEY_ATSCMH_PARADE_ID"
	case DVB_FE_KEY_ATSCMH_NOG:
		return "DVB_FE_KEY_ATSCMH_NOG"
	case DVB_FE_KEY_ATSCMH_TNOG:
		return "DVB_FE_KEY_ATSCMH_TNOG"
	case DVB_FE_KEY_ATSCMH_SGN:
		return "DVB_FE_KEY_ATSCMH_SGN"
	case DVB_FE_KEY_ATSCMH_PRC:
		return "DVB_FE_KEY_ATSCMH_PRC"
	case DVB_FE_KEY_ATSCMH_RS_FRAME_MODE:
		return "DVB_FE_KEY_ATSCMH_RS_FRAME_MODE"
	case DVB_FE_KEY_ATSCMH_RS_FRAME_ENSEMBLE:
		return "DVB_FE_KEY_ATSCMH_RS_FRAME_ENSEMBLE"
	case DVB_FE_KEY_ATSCMH_RS_CODE_MODE_PRI:
		return "DVB_FE_KEY_ATSCMH_RS_CODE_MODE_PRI"
	case DVB_FE_KEY_ATSCMH_RS_CODE_MODE_SEC:
		return "DVB_FE_KEY_ATSCMH_RS_CODE_MODE_SEC"
	case DVB_FE_KEY_ATSCMH_SCCC_BLOCK_MODE:
		return "DVB_FE_KEY_ATSCMH_SCCC_BLOCK_MODE"
	case DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_A:
		return "DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_A"
	case DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_B:
		return "DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_B"
	case DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_C:
		return "DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_C"
	case DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_D:
		return "DVB_FE_KEY_ATSCMH_SCCC_CODE_MODE_D"
	case DVB_FE_KEY_INTERLEAVING:
		return "DVB_FE_KEY_INTERLEAVING"
	case DVB_FE_KEY_LNA:
		return "DVB_FE_KEY_LNA"
	case DVB_FE_STAT_SIGNAL_STRENGTH:
		return "DVB_FE_STAT_SIGNAL_STRENGTH"
	case DVB_FE_STAT_CNR:
		return "DVB_FE_STAT_CNR"
	case DVB_FE_STAT_PRE_ERROR_BIT_COUNT:
		return "DVB_FE_STAT_PRE_ERROR_BIT_COUNT"
	case DVB_FE_STAT_PRE_TOTAL_BIT_COUNT:
		return "DVB_FE_STAT_PRE_TOTAL_BIT_COUNT"
	case DVB_FE_STAT_POST_ERROR_BIT_COUNT:
		return "DVB_FE_STAT_POST_ERROR_BIT_COUNT"
	case DVB_FE_STAT_POST_TOTAL_BIT_COUNT:
		return "DVB_FE_STAT_POST_TOTAL_BIT_COUNT"
	case DVB_FE_STAT_ERROR_BLOCK_COUNT:
		return "DVB_FE_STAT_ERROR_BLOCK_COUNT"
	case DVB_FE_STAT_TOTAL_BLOCK_COUNT:
		return "DVB_FE_STAT_TOTAL_BLOCK_COUNT"
	case DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX:
		return "DVB_FE_KEY_SCRAMBLING_SEQUENCE_INDEX"
	default:
		return "[?? Invalid DVBFrontendKey value]"
	}
}
