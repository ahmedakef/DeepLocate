package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var destination = flag.String("d", ".", "the search directory")

func main() {
	flag.Parse()
	root := *destination

	var files []FileInfo
	files = WalkSearch(root)

	b, _ := json.MarshalIndent(files, "", "\t")
	_ = ioutil.WriteFile("explored_files.json", b, 0644)

}
