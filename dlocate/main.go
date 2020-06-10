package main

import (
	"flag"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var operation = flag.String("o", "index", "the operation to do (index or search or clear")
var destination = flag.String("d", "./", "the search directory")
var searchWord = flag.String("s", "", "the search word")

var indexInfo IndexInfo
var directoryPartition DirectoryPartition

func main() {

	log.SetLevel(log.DebugLevel)

	flag.Parse()
	root := *destination
	op := *operation

	// root = "/home/ahmed/Downloads/cloud computing/"
	// op = "update"
	// remove trailling backslash
	if filepath.ToSlash(root)[len(root)-1] == '/' {
		root = root[:len(root)-1]
	}

	indexInfo = getIndexInfo()
	directoryPartition = getDirectoryPartition()

	if op == "index" {
		startIndexing(root)
	} else if op == "clear" {
		clearIndex()
	} else if op == "update" {
		update(root)
	} else if op == "searchNames" {
		word := *searchWord
		//word = "run"
		find(word, root, false)
	} else if op == "searchContent" {
		word := *searchWord
		//word = "run"
		find(word, root, true)
	} else {
		log.Info("Please select correct operation")
	}

}
