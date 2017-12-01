package main

import (
	"fmt"
	"strconv"

	"github.com/Eraac/go-gtfs"
	"github.com/train-cat/schedule-updater/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := model.Truncate()
	exitOnError(err)

	g, err := gtfs.Load("routes", nil)
	exitOnError(err)

	fs := []func(*gtfs.GTFS) error{persistCalendars, persistTrips, persistStopTimes}

	for _, f := range fs {
		exitOnError(f(g))
	}
}

func persistCalendars(g *gtfs.GTFS) error {
	log.Info("persist calendars")

	for _, cal := range g.Calendars {
		start, _ := strconv.Atoi(cal.Start)
		end, _ := strconv.Atoi(cal.End)
		days, _ := strconv.ParseInt(fmt.Sprintf("%d%d%d%d%d%d%d", cal.Monday, cal.Tuesday, cal.Wednesday, cal.Thursday, cal.Friday, cal.Saturday, cal.Sunday), 2, 8)

		c := &model.Calendar{
			CalendarID: cal.ServiceID,
			Start:      start,
			End:        end,
			Days:       int(days),
		}

		c.Persist()
	}

	return nil
}

func persistTrips(g *gtfs.GTFS) error {
	log.Info("persist trips")

	routes := getTransilienRoutes(g)

	for _, trip := range g.Trips {
		calendarID, ok := model.Calendars[trip.ServiceID]

		if !ok || !inArray(trip.RouteID, routes) {
			continue
		}

		t := &model.Trip{
			Code:       getCode(trip.ID),
			Mission:    trip.Headsign,
			CalendarID: calendarID,
			TripID:     trip.ID,
		}

		t.Persist()
	}

	return nil
}

func persistStopTimes(g *gtfs.GTFS) error {
	log.Info("persist stop_times")

	for _, stop := range g.StopsTimes {
		stationID, sok := model.Stations[getUIC(stop.StopID)]
		tripID, tok := model.Trips[stop.TripID]

		if !sok || !tok {
			continue
		}

		s := &model.StopTime{
			Schedule:  stop.Departure,
			StationID: stationID,
			TripID:    tripID,
		}

		s.Persist()
	}

	return nil
}
