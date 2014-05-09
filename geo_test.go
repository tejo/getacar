package main

import "testing"

func TestCalculateDistance(t *testing.T) {
	dist := CalculateDistance(10.0, 10.0, 20.0, 20.0)
	if dist != 1544.75756102961 {
		t.Error("error calculate distance")
	}
}

func TestClosestCars(t *testing.T) {
	setUpTestUrls()
	cs := NewCarStore()
	cs.FetchCars()
	closestCars := cs.ClosestCars(45.47665, 9.22389, "all", 100)
	if closestCars[0].Name != "118/ER805NP" {
		t.Error("error fetching closest cars")
	}
}

func TestClosestCarsWithLimit(t *testing.T) {
	limit := 10
	setUpTestUrls()
	cs := NewCarStore()
	cs.FetchCars()
	closestCars := cs.ClosestCars(45.47665, 9.22389, "all", limit)
	if len(closestCars) != limit {
		t.Error("error fetching and limiting closest cars")
	}
}

func TestClosestCarsWithFilter(t *testing.T) {
	limit := 30
	setUpTestUrls()
	cs := NewCarStore()
	cs.FetchCars()
	closestCars := cs.ClosestCars(45.47665, 9.22389, "car2go", limit)
	if cars := len(closestCars); cars != limit {
		t.Errorf("error filtering cars, len is %d, but limit was %d", cars, limit)
	}
	if !checkVendor("car2go", closestCars) {
		t.Error("mixed vendor found")
	}
}
