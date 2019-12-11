package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func visit(files *[]FileInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fileName := info.Name()
		modTime := info.ModTime()
		if !info.IsDir() {
			fileinfo := FileInfo{fileName, modTime}

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

// WalkSearch search function using Walk utility.
func WalkSearch(root string) []FileInfo {

	var files []FileInfo
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	return files
}
