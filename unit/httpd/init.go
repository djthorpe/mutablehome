/*
	Mutablehome Automation: Web Server
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package httpd

import (
	"net"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

func init() {
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     Httpd{}.Name(),
		Requires: []string{"gopi/mdns/register"},
		Config: func(app gopi.App) error {
			app.Flags().FlagString("httpd.iface", "", "Bind interface")
			app.Flags().FlagUint("httpd.port", 0, "Bind port")
			return nil
		},
		New: func(app gopi.App) (gopi.Unit, error) {
			if iface, err := interfaceForName(app.Flags().GetString("httpd.iface", gopi.FLAG_NS_DEFAULT)); err != nil {
				return nil, err
			} else {
				return gopi.New(Httpd{
					Iface: iface,
					Port:  app.Flags().GetUint("httpd.port", gopi.FLAG_NS_DEFAULT),
				}, app.Log().Clone(Httpd{}.Name()))
			}
		},
	})
}

func interfaceForName(name string) (net.Interface, error) {
	if name == "" {
		return net.Interface{}, nil
	}
	if ifaces, err := net.Interfaces(); err != nil {
		return net.Interface{}, err
	} else {
		for _, iface := range ifaces {
			if iface.Name == name {
				return iface, nil
			}
		}
	}

	// No interface found
	return net.Interface{}, gopi.ErrNotFound.WithPrefix(name)
}
