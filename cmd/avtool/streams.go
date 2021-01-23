package main

import (
	"fmt"
	"io"
	"strings"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func StreamToRow(stream *ffmpeg.AVStream) []string {
	params := stream.CodecPar()
	name := "-"
	bitrate := "-"
	extra := make([]string, 0)
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
	if w, h := params.Width(), params.Height(); w > 0 && h > 0 {
		extra = append(extra, fmt.Sprint("frame={", w, "x", h, "}"))
	}
	if d := stream.Disposition(); d != ffmpeg.AV_DISPOSITION_NONE && d != ffmpeg.AV_DISPOSITION_DEFAULT {
		extra = append(extra, d.String())
	}
	return []string{
		fmt.Sprint(stream.Id()),
		fmt.Sprint(params.Type()),
		fmt.Sprint(name),
		fmt.Sprint(bitrate),
		strings.Join(extra, ","),
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
		table.SetHeader([]string{"Id", "Type", "Codec", "Bitrate", "Parameters"})
		for _, stream := range ctx.Streams() {
			table.Append(StreamToRow(stream))
		}
		table.Render()
	}

	return nil
}
