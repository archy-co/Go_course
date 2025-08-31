package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Route struct {
	StartStation string
	EndStation   string
	NumStops     int
	Distance     float64
}

func readRoutesFromFile(filename string) ([]Route, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var routes []Route
	for _, record := range records {
		numStops, _ := strconv.Atoi(record[2])
		distance, _ := strconv.ParseFloat(record[3], 64)
		routes = append(routes, Route{
			StartStation: record[0],
			EndStation:   record[1],
			NumStops:     numStops,
			Distance:     distance,
		})
	}
	return routes, nil
}

func sortRoutesByDistance(routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Distance < routes[j].Distance
	})
}

func countRoutesWithAvgStopLengthLessThanX(routes []Route, x float64) int {
	count := 0
	for _, route := range routes {
		if route.NumStops > 0 {
			avgStopLength := route.Distance / float64(route.NumStops)
			if avgStopLength < x {
				count++
			}
		}
	}
	return count
}

func filterRoutesByStartStation(routes []Route, startStation string) []Route {
	var filteredRoutes []Route
	for _, route := range routes {
		if route.StartStation == startStation {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	return filteredRoutes
}

func findRoutesWithMaxStops(routes []Route) []Route {
	var maxStops int
	var maxRoutes []Route
	for _, route := range routes {
		if route.NumStops > maxStops {
			maxStops = route.NumStops
			maxRoutes = []Route{route}
		} else if route.NumStops == maxStops {
			maxRoutes = append(maxRoutes, route)
		}
	}
	return maxRoutes
}

func main() {
	routes, err := readRoutesFromFile("routes.csv")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	sortRoutesByDistance(routes)
	fmt.Println("Sorted Routes by Distance:", routes)

	x := 5.0
	count := countRoutesWithAvgStopLengthLessThanX(routes, x)
	fmt.Printf("Number of routes with average stop length less than %.2f km: %d\n", x, count)

	startStation := "Station A"
	filteredRoutes := filterRoutesByStartStation(routes, startStation)
	fmt.Printf("Routes starting from %s: %v\n", startStation, filteredRoutes)

	maxStopRoutes := findRoutesWithMaxStops(routes)
	fmt.Println("Routes with maximum stops:", maxStopRoutes)
}
