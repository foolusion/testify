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
	if err != nil {
		http.Error(w, "unable to hash units", http.StatusBadRequest)
	}
	uc1 := h1.uniformChoice([]string{"red", "blue", "green"})
	uc2 := h2.uniformChoice([]string{"red", "blue", "green"})
	fmt.Fprintf(w, "exp-123: units=%v, params=%v\n", units, uc1)
	fmt.Fprintf(w, "exp-456: units=%v, params=%v\n", units, uc2)
}
