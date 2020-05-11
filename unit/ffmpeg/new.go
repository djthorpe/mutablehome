/*
	ffmpeg bindings
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package ffmpeg

import (
	// Modules
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FFmpeg struct{}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (FFmpeg) Name() string { return "ffmpeg" }

func (config FFmpeg) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(ffmpeg)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}
