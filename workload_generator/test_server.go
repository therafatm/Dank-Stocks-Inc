package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"net/http"
)

var count int
var chunkSize = 1000
var start = time.Now()
var chunkStart = time.Now()

func rootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(data)
	count += 1
	if count % chunkSize == 0 {
		totalTime := time.Since(start)
		chunkTime := time.Since(chunkStart)
		fmt.Printf("%d in %s\n", count, totalTime)
		fmt.Printf("%d in %s\n", chunkSize, chunkTime)
		chunkStart = time.Now()
	}
}

func main() {
	count = 0
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}