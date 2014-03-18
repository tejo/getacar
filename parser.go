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
		entries[index].Id = entries[index].CarPlate
		entries[index].Price = 0.25
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

func fetchCarsFromAPI(car2goMilanUrl, car2goRomeUrl, enjoyUrl string) {
	data1, s1 := makeHttpRequest(car2goMilanUrl)
	data2, s2 := makeHttpRequest(enjoyUrl)
	data3, s3 := makeHttpRequest(car2goRomeUrl)
	if s1 != 200 && s2 != 200 && s3 != 200 {
		return
	}
	if s1 == 200 || s2 == 200 || s3 == 200 {
		cars = make([]CarEntry, 0)
		if s1 == 200 {
			ParseCar2GoJson(data1)
		}
		if s2 == 200 {
			ParseEnjoyJson(data2)
		}
		if s3 == 200 {
			ParseCar2GoJson(data3)
		}
	}
	log.Printf("load %d cars\n", len(cars))
	return
}

/* func fetch() { */
/* 	c := make(chan Result) */
/* 	go func() { c <- Web(query) }() */
/* 	go func() { c <- Image(query) }() */
/* 	go func() { c <- Video(query) }() */

/* 	timeout := time.After(80 * time.Millisecond) */
/* 	for i := 0; i < 3; i++ { */
/* 		select { */
/* 		case result := <-c: */
/* 			results = append(results, result) */
/* 		case <-timeout: */
/* 			fmt.Println("timed out") */
/* 			return */
/* 		} */
/* 	} */
/* 	return */
/* } */
