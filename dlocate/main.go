package main

import (
	"flag"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var indexInfo IndexInfo
var directoryPartition DirectoryPartition
var filesContent map[string]map[string]float32

var deepScan bool

func main() {

	log.SetLevel(log.DebugLevel)

	var operation, destination, searchWord string

	flag.StringVar(&operation, "o", "", "the operation to do (index or search or clear")
	flag.StringVar(&destination, "d", "./", "the search directory")
	flag.StringVar(&searchWord, "s", "", "the search word")
	flag.BoolVar(&deepScan, "deepScan", false, "Use machine learning to get content of the file")

	flag.Parse()

	// destination = "/home/ahmed/Downloads/cloud computing/"
	// operation = "web"
	// searchWord = "run"

	// remove trailling backslash
	if filepath.ToSlash(destination)[len(destination)-1] == '/' {
		destination = destination[:len(destination)-1]
	}

	indexInfo = getIndexInfo()
	directoryPartition = getDirectoryPartition()

	if operation == "index" {
		startIndexing(destination)
	} else if operation == "clear" {
		clearIndex()
	} else if operation == "update" {
		update(destination)
	} else if operation == "search" {
		find(searchWord, destination, deepScan)
	} else if operation == "web" {
		startServer()
	} else {
		log.Info("Please select correct operation")
	}

}
