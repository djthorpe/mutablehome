/*
	ffmpeg bindings
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package ffmpeg

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

func init() {
	// Gateway Connector
	gopi.UnitRegister(gopi.UnitConfig{
		Name: FFmpeg{}.Name(),
		New: func(app gopi.App) (gopi.Unit, error) {
			return gopi.New(FFmpeg{}, app.Log().Clone(FFmpeg{}.Name()))
		},
	})
}
