package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"
)

var operation = flag.String("o", "search", "the operation to do (index or search")
var destination = flag.String("d", ".", "the search directory")
var searchWord = flag.String("s", ".", "the search word")

func main() {

	log.SetLevel(log.InfoLevel)

	flag.Parse()
	root := *destination
	//word := *searchWord

	//var matchedFiles []FileInfo
	//matchedFiles = find(word, root)

	//b, _ := json.MarshalIndent(matchedFiles, "", "\t")
	//_ = ioutil.WriteFile("explored_files.json", b, 0644)

	startIndexing(root)
}
