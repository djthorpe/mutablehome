package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cmd struct {
	Name        string
	Description string
	Func        CmdFunc
}

type CmdFunc func(io.Writer, []string) error

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARS

var (
	flagDebug = flag.Bool("debug", false, "Turn on additional output")
)

var (
	commands = []Cmd{
		Cmd{"codecs", "List registered codecs", Codecs},
		Cmd{"streams", "Display stream information", Streams},
	}
)

////////////////////////////////////////////////////////////////////////////////
// MAIN

func GetCommand(args []string) (Cmd, []string) {
	// Return default command if no arguments
	if len(args) == 0 {
		return commands[0], nil
	}
	// Match command by name
	for _, command := range commands {
		if command.Name == args[0] {
			return command, args[1:]
		}
	}
	// Command not found
	return Cmd{}, nil
}

func Run() error {
	flag.Usage = Usage
	flag.Parse()

	if cmd, args := GetCommand(flag.Args()); cmd.Func != nil {
		return cmd.Func(os.Stdout, args)
	} else {
		return fmt.Errorf("Invalid command")
	}
}

func Usage() {
	w := flag.CommandLine.Output()
	name := filepath.Base(flag.CommandLine.Name())
	fmt.Fprintf(w, "Usage of %v:\n", name)

	fmt.Fprintf(w, "\nCommands:\n")
	for _, cmd := range commands {
		fmt.Fprintf(w, "  %v\n      %v\n", cmd.Name, cmd.Description)
	}

	fmt.Fprintf(w, "\nFlags:\n")
	flag.PrintDefaults()
}

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
