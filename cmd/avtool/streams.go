package main

import (
	"fmt"
	"io"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func StreamToRow(stream *ffmpeg.AVStream) []string {
	params := stream.CodecPar()
	name := "-"
	bitrate := "-"
	if codec := ffmpeg.FindDecoderById(params.Id()); codec != nil {
		name = codec.Name()
	}
	if br := params.BitRate(); br > 0 {
		if br > 1000 {
			bitrate = fmt.Sprint(br/1000, " kbps")
		} else {
			bitrate = fmt.Sprint(br, " bps")
		}
	}
	return []string{
		fmt.Sprint(stream.Id()),
		fmt.Sprint(params.Type()),
		fmt.Sprint(name),
		fmt.Sprint(bitrate),
	}
}

func Streams(w io.Writer, args []string) error {
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
		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Id", "Type", "Codec", "Bitrate"})
		for _, stream := range ctx.Streams() {
			table.Append(StreamToRow(stream))
		}
		table.Render()
	}

	return nil
}
