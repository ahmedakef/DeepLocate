package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func startServer() {

	http.HandleFunc("/search", searchWeb)
	http.HandleFunc("/index", indexWeb)
	http.HandleFunc("/clear", clearWeb)
	http.HandleFunc("/update", updateWeb)

	port := "8080"

	log.Infof("Starting development server at http://127.0.0.1:%v/", port)
	http.ListenAndServe(":"+port, nil)
}

func searchWeb(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	requiredFields := []string{"q", "destination", "deepScan"}
	err := checkKeysExits(r, requiredFields)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := r.FormValue("q")
	destination := r.FormValue("destination")
	deepScanParam := r.FormValue("deepScan")

	deepScan, err = strconv.ParseBool(deepScanParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	matchedFiles, contentMatchedFiles := find(query, destination, deepScan)

	files := map[string][]string{
		"matchedFiles":        matchedFiles,
		"contentMatchedFiles": contentMatchedFiles,
	}
	respondWithJSON(w, http.StatusOK, files)

}

func indexWeb(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	requiredFields := []string{"destination", "deepScan"}
	err := checkKeysExits(r, requiredFields)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	destination := r.FormValue("destination")
	deepScanParam := r.FormValue("deepScan")

	deepScan, err = strconv.ParseBool(deepScanParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = startIndexing(destination)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := map[string]string{
		"message": "finished indexing partitions successfully",
	}
	respondWithJSON(w, http.StatusOK, message)

}

func clearWeb(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	err := clearIndex()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := map[string]string{
		"message": "Index cleared successfully",
	}
	respondWithJSON(w, http.StatusOK, message)

}

func updateWeb(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	requiredFields := []string{"destination", "deepScan"}
	err := checkKeysExits(r, requiredFields)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	destination := r.FormValue("destination")
	deepScanParam := r.FormValue("deepScan")

	deepScan, err = strconv.ParseBool(deepScanParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = update(destination)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	message := map[string]string{
		"message": "finished updateing partitions successfully",
	}
	respondWithJSON(w, http.StatusOK, message)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
