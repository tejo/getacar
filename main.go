package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

type CarEntry struct {
	Fuel    int     `json:"fuel_level"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lon"`
	Address string  `json:"address"`
	Type    string
	//enjoy only
	CarPlate string `json:"car_plate"`
	//car2go only
	Name string
	Vin  string
}

var cars = make([]CarEntry, 0)
var car2goUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=milano&oauth_consumer_key=car2gowebsite&format=json"
var enjoyUrl string = "http://enjoy.eni.com/get_vetture"

func main() {
	startClock()
	fetchCarsFromAPI(car2goUrl, enjoyUrl)
	r := mux.NewRouter()
	r.HandleFunc("/geocode", Geocode).Methods("GET")
	r.HandleFunc("/cars", LoadCars)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func Geocode(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	addr := r.URL.Query().Get("q")
	geo := geogo.NewGeocoder()
	result := geo.Geocode(addr)
	if err := encoder.Encode(&result); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
	}
}

func LoadCars(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&cars); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
	}
}

func startClock() {
	ticker := time.NewTicker(60 * 5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchCarsFromAPI(car2goUrl, enjoyUrl)
			}
		}
	}()
}
