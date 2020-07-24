package main

import (
	utils "dlocate/osutils"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func startServer() {

	http.HandleFunc("/search", searchWeb)
	http.HandleFunc("/index", indexWeb)
	http.HandleFunc("/clear", clearWeb)
	http.HandleFunc("/update", updateWeb)
	http.HandleFunc("/metaSearch", metaSearchWeb)

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

func metaSearchWeb(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var err error
	requiredFields := []string{"q", "destination", "deepScan"}
	err = checkKeysExits(r, requiredFields)
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

	start := utils.FileMetadata{}
	end := utils.FileMetadata{}
	extentionsParam := r.FormValue("extentions")

	var extentions []string
	if len(extentionsParam) > 0 {
		extentions = strings.Split(extentionsParam, ",")
	}

	var t time.Time

	if len(r.FormValue("startATime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("startATime"))
		start.ATime = t
	}
	if len(r.FormValue("startCTime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("startCTime"))
		start.CTime = t
	}
	if len(r.FormValue("startMTime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("startMTime"))
		start.MTime = t
	}

	if len(r.FormValue("endATime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("endATime"))
		end.ATime = t
	}
	if len(r.FormValue("endCTime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("endCTime"))
		end.CTime = t
	}
	if len(r.FormValue("endMTime")) > 0 {
		t, err = time.Parse(time.RFC3339, r.FormValue("endMTime"))
		end.MTime = t
	}

	if len(r.FormValue("startSize")) > 0 {
		startSize, _ := strconv.Atoi(r.FormValue("startSize"))
		start.Size = int64(startSize)
	}

	if len(r.FormValue("endSize")) > 0 {
		endSize, _ := strconv.Atoi(r.FormValue("endSize"))
		end.Size = int64(endSize)
	}

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	matchedFiles, contentMatchedFiles := metaSearch(query, destination, deepScan,
		start, end, extentions)

	files := map[string][]string{
		"matchedFiles":        matchedFiles,
		"contentMatchedFiles": contentMatchedFiles,
	}
	respondWithJSON(w, http.StatusOK, files)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
