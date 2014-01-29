package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	dist := CalculateDistance(10.0, 10.0, 20.0, 20.0)
	if dist != 1544.75756102961 {
		t.Error("error calculate distance")
	}
}


func TestClosestCars(t *testing.T) {
	cars = make([]CarEntry, 0)
	car2go := NewTestServer("fixtures/car2go.json")
	enjoy := NewTestServer("fixtures/enjoy.json")
	fetchCarsFromAPI(car2go.URL, enjoy.URL)
	closestCars := ClosestCars(45.47665, 9.22389)
	if closestCars[0].Name != "118/ER805NP" {
		t.Error("error fetching closest cars")
	}
}

func TestParseEnjoyJson(t *testing.T) {
	cars = make([]CarEntry, 0)
	data, _ := NewMockResponse("fixtures/enjoy.json")
	ParseEnjoyJson(data)
	if len(cars) < 1 {
		t.Error("error decoding enjoy json")
	}
	if cars[0].Type != "enjoy" {
		t.Error("error post processing json")
	}
}

func TestParseCar2GoJson(t *testing.T) {
	cars = make([]CarEntry, 0)
	data, _ := NewMockResponse("fixtures/car2go.json")
	ParseCar2GoJson(data)
	if len(cars) < 1 {
		t.Error("error decoding enjoy json")
	}
	if cars[0].Type != "car2go" {
		t.Error("error post processing json")
	}
}

func TestFetchCarsFromAPI(t *testing.T) {
	cars = make([]CarEntry, 0)
	car2go := NewTestServer("fixtures/car2go.json")
	enjoy := NewTestServer("fixtures/enjoy.json")
	fetchCarsFromAPI(car2go.URL, enjoy.URL)
	if len(cars) != 689 {
		t.Error("error fetching from apis")
	}
}

/* enjoy */
/* { */
/* car_name: "Fiat 500", */
/* car_plate: "ET981CW", */
/* fuel_level: 93, */
/* lat: 45.50299, */
/* lon: 9.244517, */
/* address: "Via Palmanova, 139-143, 20132 Milano", */
/* virtual_rental_type_id: 2, */
/* virtual_rental_id: 3663, */
/* car_category_type_id: 1, */
/* car_category_id: 8, */
/* onClick_disabled: false */
/* }, */

/* car2go */
/* placemarks: [ */
/* { */
/* address: "Via Don Francesco B. della Torre, 20157 Milano", */
/* coordinates: [ */
/* 9.13938, */
/* 45.50871, */
/* 0 */
/* ], */
/* engineType: "CE", */
/* exterior: "GOOD", */
/* fuel: 54, */
/* interior: "GOOD", */
/* name: "385/ES400JF", */
/* vin: "WME4513341K647881" */
/* },turn */

func NewTestServer(fixture string) *httptest.Server {
	data, _ := NewMockResponse(fixture)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(data))
	}))
}

func NewMockResponse(s string) ([]byte, error) {
	dataPath := path.Join(s)
	_, readErr := os.Stat(dataPath)
	if readErr != nil && os.IsNotExist(readErr) {
		return nil, readErr
	}
	handler, handlerErr := os.Open(dataPath)
	if handlerErr != nil {
		return nil, handlerErr
	}

	data, readErr := ioutil.ReadAll(handler)

	if readErr != nil {
		return nil, readErr
	}

	return data, nil
}
