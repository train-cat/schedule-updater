package model

import "github.com/jinzhu/gorm"

// Station sncf
type Station struct {
	gorm.Model
	UIC string `gorm:"column:UIC"`
}

// Stations map uic <-> id
var Stations = map[string]uint{}

// LoadAllStations in memory
func LoadAllStations() error {
	var ss []Station

	err := db.Model(Station{}).Find(&ss).Error

	if err != nil {
		return err
	}

	for _, s := range ss {
		Stations[s.UIC[:len(s.UIC)-1]] = s.ID
	}

	return nil
}
