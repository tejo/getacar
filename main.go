package main

import (
	"encoding/json"
	_ "expvar"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"text/template"
	"time"

	"getacar/car2go"

	"bitbucket.org/kardianos/osext"
	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
)

var (
	cs             *CarStore
	car2goMilanUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=milano&oauth_consumer_key=getacar&format=json"
	car2goRomeUrl  string = "https://www.car2go.com/api/v2.1/vehicles?loc=roma&oauth_consumer_key=getacar&format=json"
	enjoyUrl       string = "http://enjoy.eni.com/get_vetture"
	twistUrl       string = "http://twistcar.it/assets/js/main.js"
	assetVersion   string
	homeTpl        *template.Template
)

func loadTemplate(folderPath string) {
	homeTpl = template.Must(template.New("index.html").Delims("<%", "%>").ParseFiles(folderPath + "public/index.html"))
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	homeTpl.Execute(w, map[string]interface{}{"AssetsPath": os.Getenv("CDN_PATH") + assetVersion})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	client := car2go.NewOAuthClient([]byte(os.Getenv("HASH_KEY")), []byte(os.Getenv("BLOCK_KEY")), "getacar", os.Getenv("CONSUMER_SECRET"))
	client.SetDebug(true)
	client.SetTestMode("0")

	folderPath, _ := osext.ExecutableFolder()
	loadTemplate(folderPath)
	assetVersion = fmt.Sprintf("%d", time.Now().UnixNano())

	cs = NewCarStore()
	startClock()
	cs.FetchCars()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/geocode", Geocode).Methods("GET")
	r.HandleFunc("/cars", LoadCars)
	r.HandleFunc("/cars/{lat}/{lng}", LoadClosestCars).Methods("GET")
	http.Handle("/"+assetVersion+"/", http.StripPrefix("/"+assetVersion+"/", http.FileServer(http.Dir(folderPath+"public/"))))
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
	if err := encoder.Encode(cs.cars); err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
	}
}

func LoadClosestCars(w http.ResponseWriter, r *http.Request) {
	vendor := r.URL.Query().Get("v")
	if vendor == "" {
		vendor = "all"
	}
	vars := mux.Vars(r)
	latStr := vars["lat"]
	lngStr := vars["lng"]
	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	closestCars := cs.ClosestCars(lat, lng, vendor, 30)
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
				cs.FetchCars()
				/* comment.Date = time.Now().In(time.UTC).Format(time.RFC3339) */
			}
		}
	}()
}
