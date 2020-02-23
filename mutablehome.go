/*
	Mutablehome Automation
	(c) Copyright David Thorpe 2019
	All Rights Reserved
	For Licensing and Usage information, please see LICENSE file
*/

package mutablehome

import (
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// INFLUXDB

type InfluxDB interface {
	// Create a new resultset
	NewResultSet(tags map[string]string) InfluxRS

	// Write rows to the database
	Write(InfluxRS) error
}

type InfluxRS interface {
	// Remove all rows
	RemoveAll()

	// Add a new row of metrics for a measurement name
	Add(string, map[string]interface{}) error

	// Add a new row of metrics using timestamp for a measurement name
	AddTS(string, map[string]interface{}, time.Time) error
}
