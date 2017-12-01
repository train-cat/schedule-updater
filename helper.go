package main

import (
	"strings"

	"github.com/Eraac/go-gtfs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func inArray(search string, array []string) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}

	return false
}

// DUASN124326F01001 -> 124326
func getCode(tripID string) string {
	ss := strings.Split(tripID, "-")

	str := strings.TrimLeft(ss[0], "DUASN")

	return str[:6]
}

// StopPoint:DUA8727141 -> 8727141
func getUIC(stopID string) string {
	return strings.Split(stopID, "A")[1]
}

func getTransilienRoutes(g *gtfs.GTFS) []string {
	var routes []string

	realtimeLines := viper.GetStringSlice("realtime_lines")

	for _, r := range g.Routes {
		if r.Type != gtfs.RouteTypeRail || !inArray(r.ShortName, realtimeLines) {
			continue
		}

		routes = append(routes, r.ID)
	}

	return routes
}
