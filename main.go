package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

type CarEntry struct {
	Fuel    int     `json:"fuel_level"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
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

func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	var R float64 = 6371
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	sin_dlng_2 := math.Sin(dLng / 2)
	sin_dlat_2 := math.Sin(dLat / 2)
	a := sin_dlat_2*sin_dlat_2 + sin_dlng_2*sin_dlng_2*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
