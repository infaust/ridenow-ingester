package models

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Datastore interface {
	AllLocations() ([]*Location, error)
	RelevantLocations() ([]*Location, error)
	StoreForecast(*ForecastEntry) (bool, error)
}

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
