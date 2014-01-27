package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestMain(t *testing.T) {
}

func TestParseEnjoyJson(t *testing.T) {
	data, _ := NewMockResponse("fixtures/enjoy.json")
	results := ParseEnjoyJson(data)
	if len(results) < 1 {
		t.Error("error decoding enjoy json")
	}
	if results[0].Type != "enjoy"  {
		t.Error("error post processing json")
	}
}

func TestParseCar2GoJson(t *testing.T) {
	data, _ := NewMockResponse("fixtures/car2go.json")
	results := ParseCar2GoJson(data)
	if len(results) < 1 {
		t.Error("error decoding enjoy json")
	}
	if results[0].Type != "car2go"  {
		t.Error("error post processing json")
	}
}

/* enjoy */
/* { */
/* car_name: "Fiat 500", */
/* car_plate: "ET981CW", */
/* fuel_level: 93, */
/* lat: 45.50299, */
/* lon: 9.244517, */
/* address: "Via Palmanova, 139-143, 20132 Milano", */
/* virtual_rental_type_id: 2, */
/* virtual_rental_id: 3663, */
/* car_category_type_id: 1, */
/* car_category_id: 8, */
/* onClick_disabled: false */
/* }, */

/* car2go */
/* placemarks: [ */
/* { */
/* address: "Via Don Francesco B. della Torre, 20157 Milano", */
/* coordinates: [ */
/* 9.13938, */
/* 45.50871, */
/* 0 */
/* ], */
/* engineType: "CE", */
/* exterior: "GOOD", */
/* fuel: 54, */
/* interior: "GOOD", */
/* name: "385/ES400JF", */
/* vin: "WME4513341K647881" */
/* },turn */

func NewTestServer(data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(data))
	}))
}

func NewMockResponse(s string) ([]byte, error) {
	dataPath := path.Join(s)
	_, readErr := os.Stat(dataPath)
	if readErr != nil && os.IsNotExist(readErr) {
		return nil, readErr
	}
	handler, handlerErr := os.Open(dataPath)
	if handlerErr != nil {
		return nil, handlerErr
	}

	data, readErr := ioutil.ReadAll(handler)

	if readErr != nil {
		return nil, readErr
	}

	return data, nil
}
