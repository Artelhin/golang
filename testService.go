package main

import (
	"fmt"
	"net/http"
	"io"
	"sync"
)

var data bool
var mux sync.Mutex

func assignmentHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", 400)
		return
	}
	var body []byte = make([]byte, 6)
	_, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		http.Error(w, "Unexpected error on data reading", 400)
		return
	}
	switch string(body) {
	case "accept":
		data = true
	case "reject":
		data = false
	default:
		http.Error(w, "Bad data", 400)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()
	if r.Method != "" && r.Method != "GET" {
		http.Error(w, "Only GET allowed", 400)
		return
	}
	if data {
		fmt.Fprint(w, "accept")
	} else {
		fmt.Fprint(w, "reject")
	}
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/assignment", assignmentHandler)
	http.ListenAndServe(":8080", nil)
}
