/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package tradfri

import (
	"encoding/json"
	"os"
	"path/filepath"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type token struct {
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

func (this *token) Read(path string) error {
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

func (this *token) Write(path string) error {
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
