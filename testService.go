package main

import (
	"fmt"
	"net/http"
	"io"
	"sync"
)

var (
	data bool
	mu sync.Mutex
)

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/assignment", assignmentHandler)
	http.ListenAndServe(":8080", nil)
}

func assignmentHandler(w http.ResponseWriter, r *http.Request) {
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
	mu.Lock()
	defer mu.Unlock()
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
	if r.Method != "" && r.Method != "GET" {
		http.Error(w, "Only GET allowed", 400)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if data {
		fmt.Fprint(w, "accept")
	} else {
		fmt.Fprint(w, "reject")
	}
}
