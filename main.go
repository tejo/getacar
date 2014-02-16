package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

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

var cars = make([]CarEntry, 0)
var car2goUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=milano&oauth_consumer_key=car2gowebsite&format=json"
var enjoyUrl string = "http://enjoy.eni.com/get_vetture"

func main() {
	startClock()
	fetchCarsFromAPI(car2goUrl, enjoyUrl)
	r := mux.NewRouter()
	r.HandleFunc("/geocode", Geocode).Methods("GET")
	r.HandleFunc("/cars", LoadCars)
	r.HandleFunc("/cars/{lat}/{lng}", LoadClosestCars).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
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

func LoadClosestCars(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	latStr := vars["lat"]
	lngStr := vars["lng"]
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	closestCars := ClosestCars(lat, lng)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&closestCars); err != nil {
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
				/* comment.Date = time.Now().In(time.UTC).Format(time.RFC3339) */
			}
		}
	}()
}
