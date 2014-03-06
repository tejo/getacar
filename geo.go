package main

import (
	"math"
	"sort"
)

func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	var R float64 = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	sinDlng2 := math.Sin(dLng / 2)
	sinDlat2 := math.Sin(dLat / 2)
	a := sinDlat2*sinDlat2 + sinDlng2*sinDlng2*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func ClosestCars(lat, lng float64, limit int) (result []CarEntry) {
	distances := make(map[float64]int)
	carsIndex := make([]float64, 0)
	for index, _ := range cars {
		distance := CalculateDistance(lat, lng, cars[index].Lat, cars[index].Lng)
		distances[distance] = index
		carsIndex = append(carsIndex, distance)
	}
	sort.Float64s(carsIndex)
	for index, distance := range carsIndex {
		if index < limit {
			result = append(result, cars[distances[distance]])
		} else {
			return
		}
	}
	return
}
