package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"time"
)

func setUpTestUrls() {
	enjoyUrl = NewTestServer("fixtures/enjoy.json").URL
	car2goMilanUrl = NewTestServer("fixtures/car2go.json").URL
	car2goRomeUrl = NewTestServer("fixtures/car2go.json").URL
	twistUrl = NewTestServer("fixtures/twist.js").URL
}

func checkVendor(vendor string, cars []CarEntry) bool {
	for _, car := range cars {
		if car.Type != vendor {
			return false
		}
	}
	return true
}

func NewTestServer(fixture string, ints ...int) *httptest.Server {
	data, _ := NewMockResponse(fixture)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(ints) == 1 {
			w.WriteHeader(ints[0])
		}
		if len(ints) == 2 {
			time.Sleep(time.Duration(ints[1]) * time.Millisecond)
		}
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
