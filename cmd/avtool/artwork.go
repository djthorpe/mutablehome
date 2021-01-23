package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////

func Artwork(w io.Writer, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Syntax: streams <filename>")
	}

	ctx := ffmpeg.NewAVFormatContext()
	if err := ctx.OpenInput(args[0], nil); err != nil {
		return err
	} else {
		defer ctx.CloseInput()
		if _, err := ctx.FindStreamInfo(); err != nil {
			return err
		}
		for _, stream := range ctx.Streams() {
			if stream.Disposition()&ffmpeg.AV_DISPOSITION_ATTACHED_PIC == 0 {
				continue
			} else {
				artwork := stream.AttachedPicture()
				filename := fmt.Sprint(stream.Id())
				params := stream.CodecPar()
				if codec := ffmpeg.FindDecoderById(params.Id()); codec != nil {
					filename += "." + codec.Name()
				}
				if err := ioutil.WriteFile(filename, artwork.Bytes(), 0644); err != nil {
					return err
				} else {
					fmt.Fprintln(w, "Written", strconv.Quote(filename))
				}
			}
		}
	}

	return nil
}
