package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type (
	Calendar struct {
		gorm.Model
		CalendarID string `gorm:"-"`
		Start      int
		End        int
		Days       int
	}

	StopTime struct {
		gorm.Model
		Schedule  string
		StationID uint `gorm:"column:station_id"`
		TripID    uint `gorm:"column:trip_id"`
	}

	Trip struct {
		gorm.Model
		Code       string
		Mission    string
		TripID     string `gorm:"-"`
		CalendarID uint   `gorm:"column:calendar_id"`
	}
)

var (
	// Calendars map calendar_id <-> id
	Calendars = map[string]uint{}

	// Trips map trip_id <-> id
	Trips = map[string]uint{}
)

// Truncate tables that will be recreated
func Truncate() error {
	log.Warn("truncate tables")

	tables := []string{"stop_time", "trip", "calendar"}

	if err := db.Exec("SET foreign_key_checks=0").Error; err != nil {
		return err
	}

	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE %s", table)).Error

		if err != nil {
			return err
		}
	}

	return db.Exec("SET foreign_key_checks=1").Error
}

// Persist calendar
func (c *Calendar) Persist() error {
	if err := db.Save(c).Error; err != nil {
		return err
	}

	Calendars[c.CalendarID] = c.ID

	return nil
}

// Persist StopTime
func (s *StopTime) Persist() error {
	return db.Save(s).Error
}

// Persist Trip
func (t *Trip) Persist() error {
	if err := db.Save(t).Error; err != nil {
		log.Errorf("%s - %s", t.Code, t.Mission)
		return err
	}

	Trips[t.TripID] = t.ID

	return nil
}
