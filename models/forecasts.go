package models

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

type ForecastEntry struct {
	Id              string
	LocationId      int64
	WaveHeightM     float64
	SwellPeriodSecs float64
	Created         time.Time
	Modified        time.Time
	Time            time.Time // time when the forecast applies
}

func NewForecast(wh, sp float64, locId int64, tm time.Time) *ForecastEntry {
	// create id by composing location id and forecast time
	id := strings.Join([]string{strconv.FormatInt(locId, 10), tm.Format("20060102150405")}, ":")
	now := time.Now()
	forecast := &ForecastEntry{
		Id:              id,
		LocationId:      locId,
		Time:            tm,
		Created:         now,
		Modified:        now,
		WaveHeightM:     wh,
		SwellPeriodSecs: sp}
	return forecast
}

func (db *DB) StoreForecast(fc *ForecastEntry) (bool, error) {
	query := `INSERT INTO forecast as f (id, location, wave_height_m, swell_period_sec, created, modified, time)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              ON CONFLICT (id) DO UPDATE
              SET wave_height_m = EXCLUDED.wave_height_m,
                  swell_period_sec = EXCLUDED.swell_period_sec, 
                  modified = current_timestamp
              WHERE f.wave_height_m != EXCLUDED.wave_height_m OR f.swell_period_sec != EXCLUDED.swell_period_sec;`
	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	res, err := stmt.Exec(fc.Id, fc.LocationId, fc.WaveHeightM, fc.SwellPeriodSecs, fc.Created, fc.Modified, fc.Time)
	if err != nil {
		return false, err
	}
	change, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return change > 0, err
}
