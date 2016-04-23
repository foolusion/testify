package main

import (
	"fmt"
	"log"
	"net/http"
)

var config = struct {
	globalSalt    string
	saltSeperator string
}{
	globalSalt:    "testify",
	saltSeperator: ":",
}

func init() {
	http.HandleFunc("/", rootHandler)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad query params", http.StatusBadRequest)
		return
	}
	units := r.Form["unit"]
	h, err := newHashedUnit("exp-123", units)
	if err != nil {
		http.Error(w, "unable to hash units", http.StatusBadRequest)
	}
	f := h.randomFloat(0, 100)
	i := h.randomInt(0, 100)
	fmt.Fprintf(w, "hello world: units: %v: %v, %v\n", units, f, i)
}
