package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	f, err := randomFloat("exp-123", units, 0, 100)
	if err != nil {
		http.Error(w, "couldn't hash units"+strings.Join(units, ","), http.StatusInternalServerError)
		return
	}
	i, err := randomInt("exp-123", units, 0, 100)
	if err != nil {
		http.Error(w, "couldn't hash units"+strings.Join(units, ","), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "hello world: units: %v: %v, %v\n", units, f, i)
}
