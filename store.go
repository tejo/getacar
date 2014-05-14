package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CarStore struct {
	cars    []CarEntry
	mu      sync.Mutex
	vendors []Vendor
}

func NewCarStore() *CarStore {
	return &CarStore{
		vendors: []Vendor{
			Vendor{Parser: &EnjoyParser{}, URL: enjoyUrl},
			Vendor{Parser: &Car2goParser{}, URL: car2goMilanUrl},
			Vendor{Parser: &Car2goParser{}, URL: car2goRomeUrl},
			Vendor{Parser: &TwistParser{}, URL: twistUrl},
		},
	}
}

func (c *CarStore) AddCar(car CarEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cars = append(c.cars, car)
}

func (c *CarStore) AddCars(cars []CarEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cars = append(c.cars, cars...)
}

func (c *CarStore) DestroyAll() {
	c.cars = []CarEntry{}
}

func (cs *CarStore) FetchCars() {
	cs.DestroyAll()
	c := make(chan []CarEntry)
	for _, vendor := range cs.vendors {
		go func(vendor Vendor) {
			data, status := cs.FetchHttpData(vendor.URL)
			if status == 200 {
				c <- vendor.ParseJSON(data)
			} else {
				c <- []CarEntry{}
			}
		}(vendor)
	}
	timeout := time.After(3000 * time.Millisecond)
	for _, vendor := range cs.vendors {
		select {
		case cars := <-c:
			cs.AddCars(cars)
		case <-timeout:
			log.Printf("%#v timed out", vendor)
		}
	}
	log.Printf("load %d cars\n", len(cs.cars))
}

func (c *CarStore) FetchHttpData(url string) (b []byte, status int) {
	b, status = makeHttpRequest(url)
	return
}

type CarEntry struct {
	Id        string
	UpdatedAt string
	Fuel      int     `json:"fuel_level"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lon"`
	Address   string  `json:"address"`
	Type      string
	Price     float64
	//enjoy only
	CarPlate string `json:"car_plate"`
	//car2go only
	Name string
	Vin  string
}

type Parser interface {
	ParseJSON(b []byte) (entries []CarEntry)
}

type Vendor struct {
	Parser
	URL string
}

type EnjoyParser struct{}

func (e *EnjoyParser) ParseJSON(b []byte) (entries []CarEntry) {
	json.Unmarshal(b, &entries)
	for index, _ := range entries {
		entries[index].Type = "enjoy"
		entries[index].Id = entries[index].CarPlate
		entries[index].Price = 0.25
	}
	return
}

type Car2goParser struct{}

func (c *Car2goParser) ParseJSON(b []byte) (entries []CarEntry) {
	results := make(map[string][]map[string]interface{}, 0)
	json.Unmarshal(b, &results)
	for _, car := range results["placemarks"] {
		entries = append(entries, CarEntry{
			Type:    "car2go",
			Id:      car["vin"].(string),
			Fuel:    int(car["fuel"].(float64)),
			Lat:     car["coordinates"].([]interface{})[1].(float64),
			Lng:     car["coordinates"].([]interface{})[0].(float64),
			Address: car["address"].(string),
			Name:    car["name"].(string),
			Vin:     car["vin"].(string),
			Price:   0.29,
		})
	}
	return
}

type TwistParser struct{}

func (t *TwistParser) ParseJSON(b []byte) (entries []CarEntry) {
	r := regexp.MustCompilePOSIX(`([0-9]*\.[0-9]+*\,[0-9]*\.[0-9]+)`)
	var coords []string
	for _, coordsPair := range r.FindAllString(string(b), -1) {
		coords = strings.Split(coordsPair, ",")
		lat, _ := strconv.ParseFloat(coords[0], 64)
		lng, _ := strconv.ParseFloat(coords[1], 64)
		entries = append(entries, CarEntry{
			Type:  "twist",
			Lat:   lat,
			Lng:   lng,
			Price: 0.29,
		})
	}
	return
}
