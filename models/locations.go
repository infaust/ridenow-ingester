package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type GeoPoint struct {
	Lat  float64
	Long float64
}

type Location struct {
	Id          int64
	Name        string
	Description string
	Relevant    bool
	Coords      GeoPoint
	Url         string
	Country     string
	Created     time.Time
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf("(%v %v)", p.Lat, p.Long)
}

// Value implements the driver Valuer interface and will return the string representation of the GeoPoint struct by calling the String() method
func (p *GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

// Scan implements the Scanner interface and will scan the point(x y) into the GeoPoint struct
func (p *GeoPoint) Scan(val interface{}) error {
	var source []byte
	switch val.(type) {
	case []byte:
		source = val.([]byte)
	default:
		return errors.New("Unable to perform geopoint conversion")
	}
	raw := string(source)
	vals := strings.Split(strings.Trim(raw, "()"), ",")
	lat, err := strconv.ParseFloat(vals[0], 64)
	if err != nil {
		return err
	}
	long, err := strconv.ParseFloat(vals[1], 64)
	if err != nil {
		return err
	}
	p.Lat = lat
	p.Long = long
	return nil
}

func (db *DB) parseLocations(rows *sql.Rows) ([]*Location, error) {
	locs := make([]*Location, 0)
	for rows.Next() {
		loc := new(Location)
		err := rows.Scan(&loc.Id, &loc.Name, &loc.Description,
			&loc.Relevant, &loc.Coords, &loc.Url, &loc.Country, &loc.Created)
		if err != nil {
			return nil, err
		}
		locs = append(locs, loc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locs, nil
}

func (db *DB) AllLocations() ([]*Location, error) {
	rows, err := db.Query("SELECT * FROM location")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.parseLocations(rows)
}

func (db *DB) RelevantLocations() ([]*Location, error) {
	rows, err := db.Query("SELECT * FROM location WHERE relevant")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.parseLocations(rows)
}
