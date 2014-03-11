package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

var cars = make([]CarEntry, 0)
var car2goUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=milano&oauth_consumer_key=car2gowebsite&format=json"
var enjoyUrl string = "http://enjoy.eni.com/get_vetture"

var homeTpl *template.Template = template.Must(template.New("index.html").Delims("<%", "%>").Funcs(funcs).ParseFiles("./public/index.html"))
var funcs = template.FuncMap{
	"isIt": isIt,
}

func isIt(l string) bool { return l == "it" }

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	homeTpl.Execute(w, map[string]interface{}{"CDNPath":os.Getenv("CDN_PATH")})
}

func main() {
	startClock()
	fetchCarsFromAPI(car2goUrl, enjoyUrl)
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
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
	closestCars := ClosestCars(lat, lng, 30)
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
