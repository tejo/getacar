package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ParseEnjoyJson(b []byte) {
	entries := make([]CarEntry, 0)
	json.Unmarshal(b, &entries)
	for index, _ := range entries {
		entries[index].Type = "enjoy"
		cars = append(cars, entries[index])
	}
	return
}

func ParseCar2GoJson(b []byte) {
	results := make(map[string][]map[string]interface{}, 0)
	json.Unmarshal(b, &results)
	for _, car := range results["placemarks"] {
		cars = append(cars, CarEntry{
			Type:    "car2go",
			Fuel:    int(car["fuel"].(float64)),
			Lat:     car["coordinates"].([]interface{})[1].(float64),
			Lng:     car["coordinates"].([]interface{})[0].(float64),
			Address: car["address"].(string),
			Name:    car["name"].(string),
			Vin:     car["vin"].(string),
		})
	}
	return
}

func makeHttpRequest(addr string) ([]byte, int) {
	r, err := http.Get(addr)
	if err != nil {
		log.Println(err)
		return make([]byte, 0), 500
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return make([]byte, 0), 500
	}
	return body, r.StatusCode
}

func fetchCarsFromAPI(car2goUrl, enjoyUrl string) {
	data1, s1 := makeHttpRequest(car2goUrl)
	data2, s2 := makeHttpRequest(enjoyUrl)
	if s1 != 200 && s2 != 200 {
		return
	}
	if s1 == 200 || s2 == 200 {
		cars = make([]CarEntry, 0)
		if s1 == 200 {
			ParseCar2GoJson(data1)
		}
		if s2 == 200 {
			ParseEnjoyJson(data2)
		}
	}
	log.Printf("load %d cars\n", len(cars))
	return
}

/* c := make(chan Result) */
/* go func() { c <- Web(query) } () */
/* go func() { c <- Image(query) } () */
/* go func() { c <- Video(query) } () */

/* timeout := time.After(80 * time.Millisecond) */
/*   for i := 0; i < 3; i++ { */
/*     select { */
/*       case result := <-c: */
/*                    results = append(results, result) */
/*       case <-timeout: */
/*                      fmt.Println("timed out") */
/*                        return */
/*     } */
/*   } */
/* return */
