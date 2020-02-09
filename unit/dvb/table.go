/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Table struct {
	Path string
}

type table struct {
	path     string
	sections []*section

	base.Unit
}

type section struct {
	name     string
	KeyValue map[string]string
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	reComment  = regexp.MustCompile("^[;#]\\s*(.*)$")
	reSection  = regexp.MustCompile("^\\[([^\\\\]+)\\]$")
	reKeyValue = regexp.MustCompile("^(\\w+)\\s*=\\s*(\\S*)$")
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Table) Name() string { return "mutablehome/dvb/table" }

func (config Table) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(table)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.DVBScanTable

func (this *table) Init(config Table) error {
	if config.Path == "" {
		return gopi.ErrBadParameter.WithPrefix("-dvb.path")
	} else if stat, err := os.Stat(config.Path); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(config.Path)
	} else if stat.Mode().IsRegular() == false {
		return gopi.ErrBadParameter.WithPrefix(config.Path)
	} else if fh, err := os.Open(config.Path); err != nil {
		return err
	} else {
		defer fh.Close()
		if sections, err := this.Decode(fh); err != nil {
			return err
		} else {
			this.path = config.Path
			this.sections = sections
		}
	}

	return nil
}

func (this *table) Decode(r io.Reader) ([]*section, error) {
	scanner := bufio.NewScanner(r)
	linenum := 0
	sections := make([]*section, 0, 10)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		linenum = linenum + 1
		if line == "" {
			continue
		} else if comment := reComment.FindStringSubmatch(line); len(comment) > 1 {
			// Comment - ignore
		} else if keyvalue := reKeyValue.FindStringSubmatch(line); len(keyvalue) > 1 {
			if err := this.setSectionKeyValue(sections, keyvalue[1], keyvalue[2]); err != nil {
				return nil, fmt.Errorf("Line %v: %w", linenum, err)
			}
		} else if section := reSection.FindStringSubmatch(line); len(section) > 1 {
			var err error
			if sections, err = this.setSectionName(sections, section[1]); err != nil {
				return nil, fmt.Errorf("Line %v: %w", linenum, err)
			}
		} else {
			return nil, fmt.Errorf("Line %v: Syntax Error", linenum)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return sections, nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *table) Properties() []home.DVBProperties {
	props := make([]home.DVBProperties, len(this.sections))
	for i, section := range this.sections {
		props[i] = section
	}
	return props
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *table) String() string {
	return "<" + this.Log.Name() +
		" path=" + strconv.Quote(this.path) +
		" sections=" + fmt.Sprint(this.sections) +
		">"
}

func (this *section) String() string {
	str := "<section" +
		" name=" + strconv.Quote(this.name)
	if sys, err := this.DeliverySystem(); err == nil {
		str += " delivery_system=" + fmt.Sprint(sys)
	}
	if f := this.Frequency(); f != 0 {
		str += " frequency=" + fmt.Sprint(f)
	}
	if bw := this.Bandwidth(); bw != 0 {
		str += " bandwidth=" + fmt.Sprint(bw)
	}
	if codeRateLP, err := this.CodeRateLP(); err == nil {
		str += " code_rate_lp=" + fmt.Sprint(codeRateLP)
	}
	if codeRateHP, err := this.CodeRateHP(); err == nil {
		str += " code_rate_hp=" + fmt.Sprint(codeRateHP)
	}
	if guardInterval, err := this.GuardInterval(); err == nil {
		str += " guard_interval=" + fmt.Sprint(guardInterval)
	}
	if hierarchy, err := this.Hierarchy(); err == nil {
		str += " hierarchy=" + fmt.Sprint(hierarchy)
	}
	if inversion, err := this.Inversion(); err == nil {
		str += " inversion=" + fmt.Sprint(inversion)
	}
	if modulation, err := this.Modulation(); err == nil {
		str += " modulation=" + fmt.Sprint(modulation)
	}
	if transmitMode, err := this.TransmitMode(); err == nil {
		str += " transmit_mode=" + fmt.Sprint(transmitMode)
	}
	str += " values=" + fmt.Sprint(this.KeyValue)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// SECTION

func (this *table) setSectionName(sections []*section, name string) ([]*section, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, gopi.ErrBadParameter
	}

	// Create a new section
	sections = append(sections, &section{name, make(map[string]string)})

	// Return success
	return sections, nil
}

func (this *table) setSectionKeyValue(sections []*section, key, value string) error {
	if len(sections) == 0 {
		return gopi.ErrBadParameter
	}
	if key = strings.ToUpper(key); key == "" {
		return gopi.ErrBadParameter
	}

	// Set key/value
	sections[len(sections)-1].KeyValue[key] = value

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RETURN PROPERTIES

func (this *section) Name() string {
	return this.name
}

func (this *section) DeliverySystem() (home.DVBDeliverySystem, error) {
	if value, exists := this.KeyValue["DELIVERY_SYSTEM"]; exists {
		value = MangleValue(value)
		for v := home.DVB_SYS_MIN; v <= home.DVB_SYS_MAX; v++ {
			if MatchesValue(value, v, "DVB_SYS_") {
				return v, nil
			}
		}
	}
	// Not found
	return home.DVB_SYS_NONE, gopi.ErrNotFound.WithPrefix("DELIVERY_SYSTEM")
}

func (this *section) Frequency() uint32 {
	if value, exists := this.KeyValue["FREQUENCY"]; exists {
		if value_, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(value_)
		}
	}
	// Bad parameter
	return 0
}

func (this *section) Bandwidth() uint32 {
	if value, exists := this.KeyValue["BANDWIDTH_HZ"]; exists {
		if value_, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint32(value_)
		}
	}
	// Bad parameter
	return 0
}

func (this *section) CodeRateHP() (home.DVBCodeRate, error) {
	return this.CodeRate("CODE_RATE_HP")
}

func (this *section) CodeRateLP() (home.DVBCodeRate, error) {
	return this.CodeRate("CODE_RATE_LP")
}

func (this *section) CodeRate(key string) (home.DVBCodeRate, error) {
	if value, exists := this.KeyValue[key]; exists {
		value = MangleValue(value)
		for v := home.DVB_FEC_MIN; v <= home.DVB_FEC_MAX; v++ {
			if MatchesValue(value, v, "DVB_FEC_") {
				return v, nil
			}
		}
	}
	// Not found
	return home.DVB_FEC_NONE, gopi.ErrNotFound.WithPrefix(key)
}

func (this *section) GuardInterval() (home.DVBGuardInterval, error) {
	if value, exists := this.KeyValue["GUARD_INTERVAL"]; exists {
		value = MangleValue(value)
		for v := home.DVB_GUARD_INTERVAL_MIN; v <= home.DVB_GUARD_INTERVAL_MAX; v++ {
			if MatchesValue(value, v, "DVB_GUARD_INTERVAL_") {
				return v, nil
			}
		}
	}
	// Not found
	return 0, gopi.ErrNotFound.WithPrefix("GUARD_INTERVAL")
}

func (this *section) Hierarchy() (home.DVBHierarchy, error) {
	if value, exists := this.KeyValue["HIERARCHY"]; exists {
		value = MangleValue(value)
		for v := home.DVB_HIERARCHY_MIN; v <= home.DVB_HIERARCHY_MAX; v++ {
			if MatchesValue(value, v, "DVB_HIERARCHY_") {
				return v, nil
			}
		}
	}
	// Not found
	return 0, gopi.ErrNotFound.WithPrefix("HIERARCHY")
}

func (this *section) Inversion() (home.DVBInversion, error) {
	if value, exists := this.KeyValue["INVERSION"]; exists {
		value = MangleValue(value)
		for v := home.DVB_INVERSION_MIN; v <= home.DVB_INVERSION_MAX; v++ {
			if MatchesValue(value, v, "DVB_INVERSION_") {
				return v, nil
			}
		}
	}
	// Not found
	return 0, gopi.ErrNotFound.WithPrefix("INVERSION")
}

func (this *section) Modulation() (home.DVBModulation, error) {
	if value, exists := this.KeyValue["MODULATION"]; exists {
		value = MangleValue(value)
		for v := home.DVB_MODULATION_MIN; v <= home.DVB_MODULATION_MAX; v++ {
			if MatchesValue(value, v, "DVB_MODULATION_") {
				return v, nil
			}
		}
	}
	// Not found
	return 0, gopi.ErrNotFound.WithPrefix("MODULATION")
}

func (this *section) TransmitMode() (home.DVBTransmitMode, error) {
	if value, exists := this.KeyValue["TRANSMISSION_MODE"]; exists {
		value = MangleValue(value)
		for v := home.DVB_TRANSMIT_MODE_MIN; v <= home.DVB_TRANSMIT_MODE_MAX; v++ {
			if MatchesValue(value, v, "DVB_TRANSMIT_MODE_") {
				return v, nil
			}
		}
	}
	// Not found
	return 0, gopi.ErrNotFound.WithPrefix("TRANSMISSION_MODE")
}

////////////////////////////////////////////////////////////////////////////////
// MANGLE STRING VALUE

func MangleValue(value string) string {
	parts := strings.FieldsFunc(value, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	return strings.Join(parts, "_")
}

func MatchesValue(input string, value interface{}, prefix string) bool {
	return input == strings.TrimPrefix(strings.ToUpper(fmt.Sprint(value)), prefix)
}
