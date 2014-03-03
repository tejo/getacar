package main

import "sync"

type CarStore struct {
	cars []CarEntry
	mu   sync.Mutex
}

func NewCarStore() *CarStore {
	return &CarStore{}
}

func (c *CarStore) Set(car CarEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cars = append(c.cars, car)
}

type CarEntry struct {
	Id        string
	UpdatedAt string
	Fuel      int     `json:"fuel_level"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lon"`
	Address   string  `json:"address"`
	Type      string
	//enjoy only
	CarPlate string `json:"car_plate"`
	//car2go only
	Name string
	Vin  string
}
