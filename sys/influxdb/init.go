package influxdb

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name: InfluxDB{}.Name(),
		Config: func(app gopi.App) error {
			app.Flags().FlagString("influxdb.addr", "http://localhost:8086/", "Server address")
			app.Flags().FlagString("influxdb.user", "", "Username for authentication")
			app.Flags().FlagString("influxdb.password", "", "Password for authentication")
			app.Flags().FlagString("influxdb.db", "", "Database name")
			app.Flags().FlagDuration("influxdb.timeout", 0, "Timeout for influxdb writes")
			app.Flags().FlagBool("influxdb.skipverify", false, "Skip https certificate verification")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(InfluxDB{
				Addr:       app.Flags().GetString("influxdb.addr", gopi.FLAG_NS_DEFAULT),
				User:       app.Flags().GetString("influxdb.user", gopi.FLAG_NS_DEFAULT),
				Password:   app.Flags().GetString("influxdb.password", gopi.FLAG_NS_DEFAULT),
				Database:   app.Flags().GetString("influxdb.db", gopi.FLAG_NS_DEFAULT),
				Timeout:    app.Flags().GetDuration("influxdb.timeout", gopi.FLAG_NS_DEFAULT),
				SkipVerify: app.Flags().GetBool("influxdb.skipverify", gopi.FLAG_NS_DEFAULT),
			}, app.Log().Clone(InfluxDB{}.Name()))
		},
	})
}
