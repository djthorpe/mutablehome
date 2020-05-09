package main

import (
	"io"
	"strings"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func CodecToRow(codec *ffmpeg.AVCodec) []string {
	t := strings.TrimPrefix(codec.Type().String(), "AVMEDIA_TYPE_")
	c := strings.ReplaceAll(codec.Capabilities().String(), "AV_CODEC_CAP_", "")
	return []string{
		codec.Name(),
		strings.ToLower(t),
		codec.Description(),
		strings.ToLower(c),
	}
}

func Codecs(w io.Writer, args []string) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Type", "Description", "Capabilities"})
	for _, codec := range ffmpeg.AllCodecs() {
		table.Append(CodecToRow(codec))
	}
	table.Render()
	return nil
}
