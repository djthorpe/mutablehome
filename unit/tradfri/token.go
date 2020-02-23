/*
	Mutablehome Automation: Ikea Tradfri
	(c) Copyright David Thorpe 2020
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Token struct {
	Id      string
	Token   string `json:"9091"`
	Version string `json:"9029"`
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FILENAME_TOKEN = "token.json"
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Token) CreatePath(path string) (string, error) {
	// If path is relative, then append user's home folder
	if filepath.IsAbs(path) == false {
		if home, err := os.UserHomeDir(); err != nil {
			return "", err
		} else {
			path = filepath.Join(home, path)
		}
	}
	// If path doesn't exist then try and create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			return path, err
		}
	}
	// Make sure path is available
	if stat, err := os.Stat(path); err != nil {
		return path, err
	} else if stat.IsDir() == false {
		return path, fmt.Errorf("%w: Not a folder: %v", gopi.ErrBadParameter, path)
	}
	// Success
	return path, nil
}

func (this *Token) Read(path string) error {
	filename := filepath.Join(path, FILENAME_TOKEN)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// When file doesn't exist then just empty out values
		this.Token = ""
		this.Version = ""
		return nil
	} else if fh, err := os.Open(filename); err != nil {
		return err
	} else {
		defer fh.Close()
		enc := json.NewDecoder(fh)
		if err := enc.Decode(this); err != nil {
			return err
		}
	}

	// Success
	return nil
}

func (this *Token) Write(path string) error {
	filename := filepath.Join(path, FILENAME_TOKEN)
	if fh, err := os.Create(filename); err != nil {
		return err
	} else {
		defer fh.Close()
		enc := json.NewEncoder(fh)
		if err := enc.Encode(this); err != nil {
			return err
		}
	}

	// Success
	return nil
}
