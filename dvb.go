/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	// Frameworks
	"context"

	gopi2 "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DVBDeliverySystem uint
	DVBGuardInterval  uint
	DVBHierarchy      uint
	DVBInterleaving   uint
	DVBTransmitMode   uint
	DVBModulation     uint
	DVBCodeRate       uint
	DVBInversion      uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type DVBFrontend interface {
	// Return name of the adaptor
	Name() string

	// Return supported delivery systems
	DeliverySystems() []DVBDeliverySystem

	// Tune with DVB properties, may timeout
	Tune(context.Context, DVBProperties) error

	// Implements gopi.Unit
	gopi2.Unit
}

type DVBTable interface {
	// Properties returns an array of DVB Properties
	// which define one or more signal sources
	Properties() []DVBProperties

	// Implements gopi.Unit
	gopi2.Unit
}

// DVBProperties are the properties used for reading from
// a multiplex
type DVBProperties interface {
	Name() string
	DeliverySystem() (DVBDeliverySystem, error)
	Frequency() uint32
	Bandwidth() uint32
	GuardInterval() (DVBGuardInterval, error)
	Hierarchy() (DVBHierarchy, error)
	Inversion() (DVBInversion, error)
	Modulation() (DVBModulation, error)
	TransmitMode() (DVBTransmitMode, error)
	CodeRateLP() (DVBCodeRate, error)
	CodeRateHP() (DVBCodeRate, error)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DVB_SYS_NONE         DVBDeliverySystem = 0
	DVB_SYS_DVBC_ANNEX_A DVBDeliverySystem = iota // Cable TV: DVB-C following ITU-T J.83 Annex A spec
	DVB_SYS_DVBC_ANNEX_B                          // Cable TV: DVB-C following ITU-T J.83 Annex B spec (ClearQAM)
	DVB_SYS_DVBT                                  // Terrestrial TV: DVB-T
	DVB_SYS_DSS                                   // Satellite TV: DSS (not fully supported)
	DVB_SYS_DVBS                                  // Satellite TV: DVB-S
	DVB_SYS_DVBS2                                 // Satellite TV: DVB-S2
	DVB_SYS_DVBH                                  // Terrestrial TV (mobile): DVB-H (standard deprecated)
	DVB_SYS_ISDBT                                 // Terrestrial TV: ISDB-T
	DVB_SYS_ISDBS                                 // Satellite TV: ISDB-S
	DVB_SYS_ISDBC                                 // Cable TV: ISDB-C (no drivers yet)
	DVB_SYS_ATSC                                  // Terrestrial TV: ATSC
	DVB_SYS_ATSCMH                                // Terrestrial TV (mobile): ATSC-M/H
	DVB_SYS_DTMB                                  // Terrestrial TV: DTMB
	DVB_SYS_CMMB                                  // Terrestrial TV (mobile): CMMB (not fully supported)
	DVB_SYS_DAB                                   // Digital audio: DAB (not fully supported)
	DVB_SYS_DVBT2                                 // Terrestrial TV: DVB-T2
	DVB_SYS_TURBO                                 // Satellite TV: DVB-S Turbo
	DVB_SYS_DVBC_ANNEX_C                          // Cable TV: DVB-C following ITU-T J.83 Annex C spec

	DVB_SYS_MIN = DVB_SYS_DVBC_ANNEX_A
	DVB_SYS_MAX = DVB_SYS_DVBC_ANNEX_C
)

const (
	DVB_GUARD_INTERVAL_1_32 DVBGuardInterval = iota
	DVB_GUARD_INTERVAL_1_16
	DVB_GUARD_INTERVAL_1_8
	DVB_GUARD_INTERVAL_1_4
	DVB_GUARD_INTERVAL_AUTO
	DVB_GUARD_INTERVAL_1_128
	DVB_GUARD_INTERVAL_19_128
	DVB_GUARD_INTERVAL_19_256
	DVB_GUARD_INTERVAL_PN420
	DVB_GUARD_INTERVAL_PN595
	DVB_GUARD_INTERVAL_PN945

	DVB_GUARD_INTERVAL_MIN = DVB_GUARD_INTERVAL_1_32
	DVB_GUARD_INTERVAL_MAX = DVB_GUARD_INTERVAL_PN945
)

const (
	DVB_HIERARCHY_NONE DVBHierarchy = iota
	DVB_HIERARCHY_1
	DVB_HIERARCHY_2
	DVB_HIERARCHY_4
	DVB_HIERARCHY_AUTO

	DVB_HIERARCHY_MIN = DVB_HIERARCHY_NONE
	DVB_HIERARCHY_MAX = DVB_HIERARCHY_AUTO
)

const (
	DVB_INTERLEAVING_NONE DVBInterleaving = iota
	DVB_INTERLEAVING_AUTO
	DVB_INTERLEAVING_240
	DVB_INTERLEAVING_72

	DVB_INTERLEAVING_MIN = DVB_INTERLEAVING_NONE
	DVB_INTERLEAVING_MAX = DVB_INTERLEAVING_72
)

const (
	DVB_TRANSMIT_MODE_2K DVBTransmitMode = iota
	DVB_TRANSMIT_MODE_8K
	DVB_TRANSMIT_MODE_AUTO
	DVB_TRANSMIT_MODE_4K
	DVB_TRANSMIT_MODE_1K
	DVB_TRANSMIT_MODE_16K
	DVB_TRANSMIT_MODE_32K
	DVB_TRANSMIT_MODE_C1
	DVB_TRANSMIT_MODE_C3780

	DVB_TRANSMIT_MODE_MIN = DVB_TRANSMIT_MODE_2K
	DVB_TRANSMIT_MODE_MAX = DVB_TRANSMIT_MODE_C3780
)

const (
	DVB_MODULATION_QPSK DVBModulation = iota
	DVB_MODULATION_QAM_16
	DVB_MODULATION_QAM_32
	DVB_MODULATION_QAM_64
	DVB_MODULATION_QAM_128
	DVB_MODULATION_QAM_256
	DVB_MODULATION_QAM_AUTO
	DVB_MODULATION_VSB_8
	DVB_MODULATION_VSB_16
	DVB_MODULATION_PSK_8
	DVB_MODULATION_APSK_16
	DVB_MODULATION_APSK_32
	DVB_MODULATION_DQPSK
	DVB_MODULATION_QAM_4_NR

	DVB_MODULATION_MIN = DVB_MODULATION_QPSK
	DVB_MODULATION_MAX = DVB_MODULATION_QAM_4_NR
)

const (
	DVB_FEC_NONE DVBCodeRate = iota
	DVB_FEC_1_2
	DVB_FEC_2_3
	DVB_FEC_3_4
	DVB_FEC_4_5
	DVB_FEC_5_6
	DVB_FEC_6_7
	DVB_FEC_7_8
	DVB_FEC_8_9
	DVB_FEC_AUTO
	DVB_FEC_3_5
	DVB_FEC_9_10
	DVB_FEC_2_5

	DVB_FEC_MIN = DVB_FEC_NONE
	DVB_FEC_MAX = DVB_FEC_2_5
)

const (
	DVB_INVERSION_OFF DVBInversion = iota
	DVB_INVERSION_ON
	DVB_INVERSION_AUTO

	DVB_INVERSION_MIN = DVB_INVERSION_OFF
	DVB_INVERSION_MAX = DVB_INVERSION_AUTO
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v DVBDeliverySystem) String() string {
	switch v {
	case DVB_SYS_NONE:
		return "DVB_SYS_NONE"
	case DVB_SYS_DVBC_ANNEX_A:
		return "DVB_SYS_DVBC_ANNEX_A"
	case DVB_SYS_DVBC_ANNEX_B:
		return "DVB_SYS_DVBC_ANNEX_B"
	case DVB_SYS_DVBT:
		return "DVB_SYS_DVBT"
	case DVB_SYS_DSS:
		return "DVB_SYS_DSS"
	case DVB_SYS_DVBS:
		return "DVB_SYS_DVBS"
	case DVB_SYS_DVBS2:
		return "DVB_SYS_DVBS2"
	case DVB_SYS_DVBH:
		return "DVB_SYS_DVBH"
	case DVB_SYS_ISDBT:
		return "DVB_SYS_ISDBT"
	case DVB_SYS_ISDBS:
		return "DVB_SYS_ISDBS"
	case DVB_SYS_ISDBC:
		return "DVB_SYS_ISDBC"
	case DVB_SYS_ATSC:
		return "DVB_SYS_ATSC"
	case DVB_SYS_ATSCMH:
		return "DVB_SYS_ATSCMH"
	case DVB_SYS_DTMB:
		return "DVB_SYS_DTMB"
	case DVB_SYS_CMMB:
		return "DVB_SYS_CMMB"
	case DVB_SYS_DAB:
		return "DVB_SYS_DAB"
	case DVB_SYS_DVBT2:
		return "DVB_SYS_DVBT2"
	case DVB_SYS_TURBO:
		return "DVB_SYS_TURBO"
	case DVB_SYS_DVBC_ANNEX_C:
		return "DVB_SYS_DVBC_ANNEX_C"
	default:
		return "[?? Invalid DVBDeliverySystem value]"
	}
}

func (r DVBCodeRate) String() string {
	switch r {
	case DVB_FEC_NONE:
		return "DVB_FEC_NONE"
	case DVB_FEC_1_2:
		return "DVB_FEC_1_2"
	case DVB_FEC_2_3:
		return "DVB_FEC_2_3"
	case DVB_FEC_3_4:
		return "DVB_FEC_3_4"
	case DVB_FEC_4_5:
		return "DVB_FEC_4_5"
	case DVB_FEC_5_6:
		return "DVB_FEC_5_6"
	case DVB_FEC_6_7:
		return "DVB_FEC_6_7"
	case DVB_FEC_7_8:
		return "DVB_FEC_7_8"
	case DVB_FEC_8_9:
		return "DVB_FEC_8_9"
	case DVB_FEC_AUTO:
		return "DVB_FEC_AUTO"
	case DVB_FEC_3_5:
		return "DVB_FEC_3_5"
	case DVB_FEC_9_10:
		return "DVB_FEC_9_10"
	case DVB_FEC_2_5:
		return "DVB_FEC_2_5"
	default:
		return "[?? Invalid DVBCodeRate value]"
	}
}

func (v DVBGuardInterval) String() string {
	switch v {
	case DVB_GUARD_INTERVAL_1_32:
		return "DVB_GUARD_INTERVAL_1_32"
	case DVB_GUARD_INTERVAL_1_16:
		return "DVB_GUARD_INTERVAL_1_16"
	case DVB_GUARD_INTERVAL_1_8:
		return "DVB_GUARD_INTERVAL_1_8"
	case DVB_GUARD_INTERVAL_1_4:
		return "DVB_GUARD_INTERVAL_1_4"
	case DVB_GUARD_INTERVAL_AUTO:
		return "DVB_GUARD_INTERVAL_AUTO"
	case DVB_GUARD_INTERVAL_1_128:
		return "DVB_GUARD_INTERVAL_1_128"
	case DVB_GUARD_INTERVAL_19_128:
		return "DVB_GUARD_INTERVAL_19_128"
	case DVB_GUARD_INTERVAL_19_256:
		return "DVB_GUARD_INTERVAL_19_256"
	case DVB_GUARD_INTERVAL_PN420:
		return "DVB_GUARD_INTERVAL_PN420"
	case DVB_GUARD_INTERVAL_PN595:
		return "DVB_GUARD_INTERVAL_PN595"
	case DVB_GUARD_INTERVAL_PN945:
		return "DVB_GUARD_INTERVAL_PN945"
	default:
		return "[?? Invalid DVBGuardInterval value]"
	}
}

func (v DVBHierarchy) String() string {
	switch v {
	case DVB_HIERARCHY_NONE:
		return "DVB_HIERARCHY_NONE"
	case DVB_HIERARCHY_1:
		return "DVB_HIERARCHY_1"
	case DVB_HIERARCHY_2:
		return "DVB_HIERARCHY_2"
	case DVB_HIERARCHY_4:
		return "DVB_HIERARCHY_4"
	case DVB_HIERARCHY_AUTO:
		return "DVB_HIERARCHY_AUTO"
	default:
		return "[?? Invalid DVBHierarchy value]"
	}
}

func (v DVBInterleaving) String() string {
	switch v {
	case DVB_INTERLEAVING_NONE:
		return "DVB_INTERLEAVING_NONE"
	case DVB_INTERLEAVING_AUTO:
		return "DVB_INTERLEAVING_AUTO"
	case DVB_INTERLEAVING_240:
		return "DVB_INTERLEAVING_240"
	case DVB_INTERLEAVING_72:
		return "DVB_INTERLEAVING_72"
	default:
		return "[?? Invalid DVBInterleaving value]"
	}
}

func (v DVBTransmitMode) String() string {
	switch v {
	case DVB_TRANSMIT_MODE_2K:
		return "DVB_TRANSMIT_MODE_2K"
	case DVB_TRANSMIT_MODE_8K:
		return "DVB_TRANSMIT_MODE_8K"
	case DVB_TRANSMIT_MODE_AUTO:
		return "DVB_TRANSMIT_MODE_AUTO"
	case DVB_TRANSMIT_MODE_4K:
		return "DVB_TRANSMIT_MODE_4K"
	case DVB_TRANSMIT_MODE_1K:
		return "DVB_TRANSMIT_MODE_1K"
	case DVB_TRANSMIT_MODE_16K:
		return "DVB_TRANSMIT_MODE_16K"
	case DVB_TRANSMIT_MODE_32K:
		return "DVB_TRANSMIT_MODE_32K"
	case DVB_TRANSMIT_MODE_C1:
		return "DVB_TRANSMIT_MODE_C1"
	case DVB_TRANSMIT_MODE_C3780:
		return "DVB_TRANSMIT_MODE_C3780"
	default:
		return "[?? Invalid DVBTransmitMode value]"
	}
}

func (v DVBModulation) String() string {
	switch v {
	case DVB_MODULATION_QPSK:
		return "DVB_MODULATION_QPSK"
	case DVB_MODULATION_QAM_16:
		return "DVB_MODULATION_QAM_16"
	case DVB_MODULATION_QAM_32:
		return "DVB_MODULATION_QAM_32"
	case DVB_MODULATION_QAM_64:
		return "DVB_MODULATION_QAM_64"
	case DVB_MODULATION_QAM_128:
		return "DVB_MODULATION_QAM_128"
	case DVB_MODULATION_QAM_256:
		return "DVB_MODULATION_QAM_256"
	case DVB_MODULATION_QAM_AUTO:
		return "DVB_MODULATION_QAM_AUTO"
	case DVB_MODULATION_VSB_8:
		return "DVB_MODULATION_VSB_8"
	case DVB_MODULATION_VSB_16:
		return "DVB_MODULATION_VSB_16"
	case DVB_MODULATION_PSK_8:
		return "DVB_MODULATION_PSK_8"
	case DVB_MODULATION_APSK_16:
		return "DVB_MODULATION_APSK_16"
	case DVB_MODULATION_APSK_32:
		return "DVB_MODULATION_APSK_32"
	case DVB_MODULATION_DQPSK:
		return "DVB_MODULATION_DQPSK"
	case DVB_MODULATION_QAM_4_NR:
		return "DVB_MODULATION_QAM_4_NR"
	default:
		return "[?? Invalid DVBModulation value]"
	}
}

func (v DVBInversion) String() string {
	switch v {
	case DVB_INVERSION_OFF:
		return "DVB_INVERSION_OFF"
	case DVB_INVERSION_ON:
		return "DVB_INVERSION_ON"
	case DVB_INVERSION_AUTO:
		return "DVB_INVERSION_AUTO"
	default:
		return "[?? Invalid DVBInversion value]"
	}
}
