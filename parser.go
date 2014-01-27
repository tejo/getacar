package main

import "encoding/json"

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

func ParseEnjoyJson(b []byte) (entries []CarEntry) {
	json.Unmarshal(b, &entries)
	for index, _ := range entries {
		entries[index].Type = "enjoy"
	}
	return
}

func ParseCar2GoJson(b []byte) (entries []CarEntry) {
	results := make(map[string][]map[string]interface{}, 0)
	json.Unmarshal(b, &results)
	for _, car := range results["placemarks"] {
		entries = append(entries, CarEntry{
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
