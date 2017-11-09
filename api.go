package main

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(response interface{}, w http.ResponseWriter) {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Super tajne!"))
}
