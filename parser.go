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
