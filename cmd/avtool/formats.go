package main

import (
	"io"
	"strings"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func OutputFormatToRow(format *ffmpeg.AVOutputFormat) []string {
	f := strings.ReplaceAll(format.Flags().String(), "AVFMT_", "")
	if f == "NONE" {
		f = ""
	}
	return []string{
		format.Name(),
		format.Description(),
		strings.ReplaceAll(format.Ext(), ",", " "),
		strings.ReplaceAll(format.MimeType(), ",", " "),
		strings.ToLower(f),
	}
}

func InputFormatToRow(format *ffmpeg.AVInputFormat) []string {
	f := strings.ReplaceAll(format.Flags().String(), "AVFMT_", "")
	if f == "NONE" {
		f = ""
	}
	return []string{
		format.Name(),
		format.Description(),
		strings.ReplaceAll(format.Ext(), ",", " "),
		strings.ReplaceAll(format.MimeType(), ",", " "),
		strings.ToLower(f),
	}
}

func Muxers(w io.Writer, args []string) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Description", "Ext", "Mimetype", "Flags"})
	for _, format := range ffmpeg.AllMuxers() {
		table.Append(OutputFormatToRow(format))
	}
	table.Render()
	return nil
}

func Demuxers(w io.Writer, args []string) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Description", "Ext", "Mimetype", "Flags"})
	for _, format := range ffmpeg.AllDemuxers() {
		table.Append(InputFormatToRow(format))
	}
	table.Render()
	return nil
}
