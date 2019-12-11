package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("you must provide search directory")
		return
	}
	root := os.Args[1]

	var files []FileInfo
	files = WalkSearch(root)

	b, _ := json.MarshalIndent(files, "", "\t")
	_ = ioutil.WriteFile("explored_files.json", b, 0644)

}
