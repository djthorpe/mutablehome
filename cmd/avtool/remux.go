package main

import (
	"fmt"
	"io"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////

func Remux(w io.Writer, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Syntax: remux <in> <out>")
	}

	// Read input
	ctxIn := ffmpeg.NewAVFormatContext()
	if err := ctxIn.OpenInput(args[0], nil); err != nil {
		return err
	}
	defer ctxIn.CloseInput()
	if _, err := ctxIn.FindStreamInfo(); err != nil {
		return err
	}

	// Determine output stream
	if ctxOut, err := ffmpeg.NewAVFormatOutputContext(args[1], nil); err != nil {
		return err
	} else {
		defer ctxOut.Free()
		streamsOut := make([]*ffmpeg.AVStream, 0, ctxIn.NumStreams())
		for _, streamIn := range ctxIn.Streams() {
			if codecInParams := streamIn.CodecPar(); codecInParams.Type() != ffmpeg.AVMEDIA_TYPE_AUDIO && codecInParams.Type() != ffmpeg.AVMEDIA_TYPE_VIDEO && codecInParams.Type() != ffmpeg.AVMEDIA_TYPE_SUBTITLE {
				continue
			}
			if streamOut := ffmpeg.NewStream(ctxOut, nil); streamOut == nil {
				return err
			} else if err := streamOut.CodecPar().From(streamIn.CodecPar()); err != nil {
				return err
			} else {
				streamsOut = append(streamsOut, streamOut)
			}
		}
		ctxOut.Dump(0)
	}

	return nil
}
