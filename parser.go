package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

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
