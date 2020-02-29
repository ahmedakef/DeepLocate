package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

var destination = flag.String("d", ".", "the search directory")
var searchWord = flag.String("s", ".", "the search word")

func main() {

	log.SetLevel(log.InfoLevel)

	flag.Parse()
	root := *destination
	word := *searchWord

	var matchedFiles []FileInfo
	//matchedFiles = find(word, root)

	b, _ := json.MarshalIndent(matchedFiles, "", "\t")
	_ = ioutil.WriteFile("explored_files.json", b, 0644)

	//indexLucene(root)
	//findLucene(word)

	indexBleve(root)
	findBleve(word)
}
