/*
	ffmpeg bindings
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package ffmpeg

import (
	"errors"

	base "github.com/djthorpe/gopi/v2/base"
	ff "github.com/djthorpe/mutablehome/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ffmpeg struct {
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func (this *ffmpeg) Init(config FFmpeg) error {
	ff.AVFormatInit()
	ff.AVDeviceInit()

	// Set logging
	level := ff.AV_LOG_INFO
	if this.Log.IsDebug() {
		level = ff.AV_LOG_DEBUG
	}
	ff.AVLogSetCallback(level, func(level ff.AVLogLevel, message string, userInfo uintptr) {
		switch level {
		case ff.AV_LOG_DEBUG, ff.AV_LOG_TRACE:
			this.Log.Debug(message)
		case ff.AV_LOG_INFO:
			this.Log.Info(message)
		case ff.AV_LOG_WARNING:
			this.Log.Warn(message)
		case ff.AV_LOG_ERROR, ff.AV_LOG_FATAL, ff.AV_LOG_PANIC:
			this.Log.Error(errors.New(message))
		}
	})

	// Return success
	return nil
}

func (this *ffmpeg) Close() error {

	// Deinit format
	ff.AVFormatDeinit()

	// Return success
	return nil
}
