package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

func main() {
	doStuff()
	r := mux.NewRouter()
	r.HandleFunc("/geocode", Geocode).Methods("GET")
	r.HandleFunc("/load/{dataId}", Load)
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

func Load(w http.ResponseWriter, r *http.Request) {

}

func doStuff() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println(time.Now())
			}
		}
	}()
}
















