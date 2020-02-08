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

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	home "github.com/djthorpe/mutablehome"
	dvb "github.com/djthorpe/mutablehome/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Table struct {
	Path string
}

type table struct {
	log      gopi.Logger
	path     string
	sections []*section

	base.Unit
}

type section struct {
	Name     string
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
		" name=" + strconv.Quote(this.Name)
	if sys := this.DeliverySystem(); sys != dvb.DVB_FE_SYS_NONE {
		str += " delivery_system=" + fmt.Sprint(this.DeliverySystem())
	}
	if f := this.Frequency(); f != 0 {
		str += " frequency=" + fmt.Sprint(f)
	}
	if bw := this.Bandwidth(); bw != 0 {
		str += " bandwidth=" + fmt.Sprint(bw)
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
	sections = append(sections, &section{
		Name:     name,
		KeyValue: make(map[string]string),
	})

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

func (this *section) DeliverySystem() dvb.DVBFEDeliverySystem {
	if value, exists := this.KeyValue["DELIVERY_SYSTEM"]; exists {
		value = strings.ToUpper(value)
		for v := dvb.DVB_FE_SYS_MIN; v <= dvb.DVB_FE_SYS_MAX; v++ {
			vstr := strings.TrimPrefix(fmt.Sprint(v), "DVB_FE_SYS_")
			if vstr == value {
				return v
			}
		}
	}

	// Not found
	return dvb.DVB_FE_SYS_NONE
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
