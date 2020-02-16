/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	home "github.com/djthorpe/mutablehome"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	frontend := app.UnitInstance("mutablehome/dvb/frontend").(home.DVBFrontend)
	demux := app.UnitInstance("mutablehome/dvb/demux").(home.DVBDemux)

	// Obtain tuning parameters
	key := app.Flags().GetString("dvb.name", gopi.FLAG_NS_DEFAULT)
	props, err := GetTransmitters(app, key)
	if err != nil {
		return err
	}

	if len(props) != 1 {
		props, _ := GetTransmitters(app, "")
		names := ""
		for _, prop := range props {
			names += strconv.Quote(prop.Name()) + ","
		}
		return fmt.Errorf("Use -dvb.name argument to select from %v", strings.TrimSuffix(names, ","))
	}

	// Tune
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	fmt.Println("Tune", props[0].Name(), "frequency=", props[0].Frequency(), "Hz")
	if err := frontend.Tune(ctx, props[0]); err != nil {
		return err
	}

	if len(args) > 0 {
		// Stream
		if pids, err := GetPids(args); err != nil {
			return err
		} else if filter, err := demux.NewStreamFilter(pids); err != nil {
			return err
		} else {
			fmt.Println("filter=", filter)
		}
	} else {
		// Initiate program association table (PAT) scanning
		if _, err := demux.ScanPAT(); err != nil {
			return err
		}
		/*
			// Initiate service description table (SDT) scanning
			if _, err := demux.ScanSDT(false); err != nil {
				return err
			}
			if _, err := demux.ScanSDT(true); err != nil {
				return err
			}
			// Initiate network information table (NIT) scanning
			if _, err := demux.ScanNIT(false); err != nil {
				return err
			}
			if _, err := demux.ScanNIT(true); err != nil {
				return err
			}
		*/ // Initiate event information (now/next) scanning
		/*
			if _, err := demux.ScanEITNowNext(false); err != nil {
				return err
			}
			if _, err := demux.ScanEITNowNext(true); err != nil {
				return err
			}
		*/
	}

	fmt.Println("Wait for CTRL+C")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetTransmitters(app gopi.App, key string) ([]home.DVBProperties, error) {
	table := app.UnitInstance("mutablehome/dvb/table").(home.DVBTable)

	if props := table.Properties(key); len(props) == 0 {
		return nil, gopi.ErrNotFound.WithPrefix("-dvb.name")
	} else {
		return props, nil
	}
}

func GetPids(args []string) ([]uint16, error) {
	pids := make([]uint16, len(args))
	for i, arg := range args {
		if pid, err := strconv.ParseUint(arg, 10, 32); err != nil {
			return nil, fmt.Errorf("Invalid pid: %v", strconv.Quote(arg))
		} else if pid > 8192 {
			return nil, fmt.Errorf("Invalid pid: %v", strconv.Quote(arg))
		} else {
			pids[i] = uint16(pid)
		}
	}
	return pids, nil
}

////////////////////////////////////////////////////////////////////////////////

func DVBSectionEventHandler(ctx context.Context, app gopi.App, evt gopi.Event) {
	demux := app.UnitInstance("mutablehome/dvb/demux").(home.DVBDemux)
	section := evt.(home.DVBSectionEvent).Section()
	filter := evt.(home.DVBSectionEvent).Filter()

	switch section.Type() {
	case home.DVB_TS_TABLE_PAT:
		// Stop Filter
		if err := demux.DestroyFilter(filter); err != nil {
			app.Log().Error(err)
		}
		// Scan for Program map specific data (PMT)
		if _, err := demux.ScanPMT(section); err != nil {
			app.Log().Error(err)
		}
	case home.DVB_TS_TABLE_PMT:
		fmt.Println(section)
		// Stop Filter
		if err := demux.DestroyFilter(filter); err != nil {
			app.Log().Error(err)
		}
	case home.DVB_TS_TABLE_SDT, home.DVB_TS_TABLE_SDT_OTHER:
		fmt.Println(section)
		// Stop Filter
		if err := demux.DestroyFilter(filter); err != nil {
			app.Log().Error(err)
		}
	case home.DVB_TS_TABLE_NIT, home.DVB_TS_TABLE_NIT_OTHER:
		fmt.Println(section)
		// Stop Filter
		if err := demux.DestroyFilter(filter); err != nil {
			app.Log().Error(err)
		}
	case home.DVB_TS_TABLE_EIT, home.DVB_TS_TABLE_EIT_OTHER:
		fmt.Println(section)
	default:
		app.Log().Warn("DVBSectionEventHandler: Unhandled:", section.Type())
	}
}
