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
	h1, err := newHashedUnit("exp-123", units)
	h2, err := newHashedUnit("exp-456", units)
	h3, err := newHashedUnit("exp-789", units)
	if err != nil {
		http.Error(w, "unable to hash units", http.StatusBadRequest)
		return
	}
	uc1 := h1.uniformChoice([]string{"red", "blue", "green"})
	uc2 := h2.uniformChoice([]string{"red", "blue", "green"})
	wc1, err := h3.weightedChoice([]string{"red", "blue", "green"}, []float64{1, 2, 1})
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "exp-123: units=%v, params=%v\n", units, uc1)
	fmt.Fprintf(w, "exp-456: units=%v, params=%v\n", units, uc2)
	fmt.Fprintf(w, "exp-789: units=%v, params=%v\n", units, wc1)
}
