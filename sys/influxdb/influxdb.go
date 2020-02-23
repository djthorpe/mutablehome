package influxdb

import (
	// Frameworks

	"time"

	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	"github.com/djthorpe/mutablehome"

	// InfluxDB Client
	_ "github.com/influxdata/influxdb1-client"
	db "github.com/influxdata/influxdb1-client/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InfluxDB struct {
	Addr       string
	User       string
	Password   string
	Timeout    time.Duration
	SkipVerify bool
	Database   string
}

type influxdb struct {
	client   db.Client
	database string

	base.Unit
}

type resultset struct {
	database string
	tags     map[string]string
	rows     db.BatchPoints
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (InfluxDB) Name() string { return "influxdb/v1" }

func (config InfluxDB) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(influxdb)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *influxdb) Init(config InfluxDB) error {
	if client, err := db.NewHTTPClient(db.HTTPConfig{
		Addr:               config.Addr,
		Username:           config.User,
		Password:           config.Password,
		Timeout:            config.Timeout,
		InsecureSkipVerify: config.SkipVerify,
	}); err != nil {
		return err
	} else {
		this.client = client
	}

	if config.Database == "" {
		return gopi.ErrBadParameter.WithPrefix("database")
	} else {
		this.database = config.Database
	}

	// Ping client
	if _, _, err := this.client.Ping(config.Timeout); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *influxdb) Close() error {

	if this.client != nil {
		this.client.Close()
	}

	// Release resources
	this.client = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION mutablehome.InfluxDB

// Create ResultSet object
func (this *influxdb) NewResultSet(tags map[string]string) mutablehome.InfluxRS {
	if len(tags) == 0 {
		return nil
	} else {
		return &resultset{this.database, tags, nil}
	}
}

// Write ResultSet to database, empty points after written without error
func (this *influxdb) Write(rs mutablehome.InfluxRS) error {
	if rs_, ok := rs.(*resultset); ok == false {
		return gopi.ErrBadParameter.WithPrefix("resultset")
	} else if rs_.rows == nil {
		return gopi.ErrNotFound.WithPrefix("resultset")
	} else if err := this.client.Write(rs_.rows); err != nil {
		return err
	} else {
		rs.RemoveAll()
		return nil
	}
}

// Empty fields from resultset (retaining tags)
func (this *resultset) RemoveAll() {
	this.rows = nil
}

// Add data using InfluxDB timestamp
func (this *resultset) Add(measurement string, fields map[string]interface{}) error {
	// Create BatchPoints
	if this.rows == nil {
		if rows, err := db.NewBatchPoints(db.BatchPointsConfig{
			Database: this.database,
		}); err != nil {
			return err
		} else {
			this.rows = rows
		}
	}

	// Append points
	if pt, err := db.NewPoint(measurement, this.tags, fields); err != nil {
		return err
	} else {
		this.rows.AddPoint(pt)
	}

	// Return success
	return nil
}

// Add data with timestamp
func (this *resultset) AddTS(measurement string, fields map[string]interface{}, timestamp time.Time) error {
	// Create BatchPoints
	if this.rows == nil {
		if rows, err := db.NewBatchPoints(db.BatchPointsConfig{
			Database: this.database,
		}); err != nil {
			return err
		} else {
			this.rows = rows
		}
	}

	// Append points
	if pt, err := db.NewPoint(measurement, this.tags, fields, timestamp); err != nil {
		return err
	} else {
		this.rows.AddPoint(pt)
	}

	// Return success
	return nil
}
