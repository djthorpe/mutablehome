/*
  Mutablehome Automation: DVB
  (c) Copyright David Thorpe 2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE file
*/

package dvb

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/mutablehome"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: Table{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagString("dvb.path", "", "DVB Multiplexer file")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Table{
				Path: app.Flags().GetString("dvb.path", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(Table{}.Name()))
		},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name: Frontend{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagUint("dvb.adapter", 0, "DVB Adapter")
			app.Flags().FlagUint("dvb.frontend", 0, "DVB Frontend")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Frontend{
				Adapter:  app.Flags().GetUint("dvb.adapter", gopi.FLAG_NS_DEFAULT),
				Frontend: app.Flags().GetUint("dvb.frontend", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(Frontend{}.Name()))
		},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Demux{}.Name(),
		Requires: []string{Frontend{}.Name(), "gopi/filepoll", "gopi/bus"},
		Config: func(app gopi.App) error {
			app.Flags().FlagUint("dvb.demux", 0, "DVB Demux")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(Demux{
				Adapter:  app.Flags().GetUint("dvb.adapter", gopi.FLAG_NS_DEFAULT),
				Demux:    app.Flags().GetUint("dvb.demux", gopi.FLAG_NS_DEFAULT),
				Frontend: app.UnitInstance(Frontend{}.Name()).(mutablehome.DVBFrontend),
				FilePoll: app.UnitInstance("gopi/filepoll").(gopi.FilePoll),
				Bus:      app.UnitInstance("gopi/bus").(gopi.Bus),
			}, app.Log().Clone(Demux{}.Name()))
		},
	})
}
