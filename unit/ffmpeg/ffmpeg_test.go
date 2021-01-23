package ffmpeg_test

import (
	"testing"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/mutablehome/unit/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////

func Test_FFmpeg_000(t *testing.T) {
	t.Log("Test_Tradfri_000")
}

////////////////////////////////////////////////////////////////////////////////

func Test_FFmpeg_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_FFmpeg_001, nil, "ffmpeg"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_FFmpeg_001(app gopi.App, t *testing.T) {
	ffmpeg := app.UnitInstance("ffmpeg")
	t.Log(ffmpeg)
}
