package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func checkKeysExits(r *http.Request, keys []string) error {
	for _, key := range keys {
		value := r.FormValue(key)
		if len(value) == 0 {
			return fmt.Errorf("You should enter a %v", key)
		}
	}
	return nil
}
