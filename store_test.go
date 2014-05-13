package main

import "testing"

func TestAddCar(t *testing.T) {
	carStore := NewCarStore()
	carStore.AddCar(CarEntry{})
	if len(carStore.cars) != 1 {
		t.Error("error adding car")
	}
}

func TestAddCars(t *testing.T) {
	cars := []CarEntry{CarEntry{}, CarEntry{}}
	carStore := NewCarStore()
	carStore.AddCars(cars)
	if len(carStore.cars) != 2 {
		t.Error("error adding cars")
	}
}

func TestFetchCars(t *testing.T) {
	setUpTestUrls()
	carStore := NewCarStore()
	carStore.FetchCars()
	carStore.FetchCars()
	if len(carStore.cars) != 1146 {
		t.Error("error fetching cars")
	}
}

func TestFetchCarsWith500Error(t *testing.T) {
	setUpTestUrls()
	twistUrl = NewTestServer("fixtures/twist.js", 500).URL
	carStore := NewCarStore()
	carStore.FetchCars()
	if len(carStore.cars) == 1146 {
		t.Error("error fetching cars")
	}
}

func TestParseEnjoyJson(t *testing.T) {
	data, _ := NewMockResponse("fixtures/enjoy.json")
	e := &Enjoy{}
	cars := e.ParseJson(data)
	if len(cars) < 1 {
		t.Error("error decoding enjoy json")
	}
	if cars[0].Type != "enjoy" {
		t.Error("error post processing json")
	}
}

func TestParseTwist(t *testing.T) {
	data, _ := NewMockResponse("fixtures/twist.js")
	tw := &Twist{}
	cars := tw.ParseJson(data)
	if len(cars) < 1 {
		t.Error("error decoding enjoy json")
	}
	if cars[0].Type != "twist" {
		t.Error("error post processing json")
	}
}

func TestParseCar2GoJson(t *testing.T) {
	data, _ := NewMockResponse("fixtures/car2go.json")
	c := &Car2go{}
	cars := c.ParseJson(data)
	if len(cars) < 1 {
		t.Error("error decoding enjoy json")
	}
	if cars[0].Type != "car2go" {
		t.Error("error post processing json")
	}
}
