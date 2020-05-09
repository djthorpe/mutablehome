package main

import (
	"fmt"
	"io"

	// Modules
	"github.com/djthorpe/mutablehome/sys/ffmpeg"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func Metadata(w io.Writer, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Syntax: metadata <filename>")
	}

	ctx := ffmpeg.NewAVFormatContext()
	if err := ctx.OpenInput(args[0], nil); err != nil {
		return err
	} else {
		defer ctx.CloseInput()
		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Name", "Value"})
		entries := ctx.Metadata().Entries()
		for _, entry := range entries {
			table.Append([]string{
				entry.Key(),
				entry.Value(),
			})
		}
		table.Render()
	}

	return nil
}
