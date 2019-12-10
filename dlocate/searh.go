package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
)

type fileInfo struct {
	FileName string
	ModTime  time.Time
}

func visit(files *[]fileInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fileName := info.Name()
		modTime := info.ModTime()
		if !info.IsDir() {
			fileinfo := fileInfo{fileName, modTime}

			log.WithFields(log.Fields{
				"fileName": fileName,
				"modTime":  modTime,
			}).Info("enter this file")
			*files = append(*files, fileinfo)
		} else {
			log.WithFields(log.Fields{
				"fileName": fileName,
				"modTime":  modTime,
			}).Info("enter this directory")
		}
		return nil
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("you must provide search directory")
		return
	}

	var files []fileInfo
	root := os.Args[1]
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(files, "", "\t")
	_ = ioutil.WriteFile("explored_files.json", b, 0644)

}
