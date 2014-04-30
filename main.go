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

	"./car2go"

	"bitbucket.org/kardianos/osext"
	"github.com/gorilla/mux"
	"github.com/tejo/geogo"
	"github.com/yvasiyarov/gorelic"
)

var cars = make([]CarEntry, 0)
var car2goMilanUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=milano&oauth_consumer_key=getacar&format=json"
var car2goRomeUrl string = "https://www.car2go.com/api/v2.1/vehicles?loc=roma&oauth_consumer_key=getacar&format=json"
var enjoyUrl string = "http://enjoy.eni.com/get_vetture"

var assetVersion string
var homeTpl *template.Template
var funcs = template.FuncMap{
	"isIt": isIt,
}

func isIt(l string) bool { return l == "it" }

func loadTemplate(folderPath string) {
	homeTpl = template.Must(template.New("index.html").Delims("<%", "%>").Funcs(funcs).ParseFiles(folderPath + "public/index.html"))
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	homeTpl.Execute(w, map[string]interface{}{"AssetsPath": os.Getenv("CDN_PATH") + assetVersion})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	agent := gorelic.NewAgent()
	agent.Verbose = false
	agent.NewrelicLicense = "ce9a50394a44f5cdc815cc318ffa0915cdeecad4"
	agent.Run()

	client := car2go.NewOAuthClient([]byte("***REMOVED***"), []byte("***REMOVED***"), "getacar", "***REMOVED***")
	client.SetDebug(true)

	folderPath, _ := osext.ExecutableFolder()
	loadTemplate(folderPath)
	assetVersion = fmt.Sprintf("%d", time.Now().UnixNano())

	startClock()
	fetchCarsFromAPI(car2goMilanUrl, car2goRomeUrl, enjoyUrl)

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
				fetchCarsFromAPI(car2goMilanUrl, car2goRomeUrl, enjoyUrl)
				/* comment.Date = time.Now().In(time.UTC).Format(time.RFC3339) */
			}
		}
	}()
}
