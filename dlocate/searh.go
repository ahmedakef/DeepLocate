package main

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func getType(info os.FileInfo) FileType {
	dot := strings.LastIndex(info.Name(), ".")
	extention := info.Name()[dot+1:]

	var fileType FileType

	if extention == "mp3" {
		fileType = "audio"
	} else if extention == "mp4" {
		fileType = "video"
	} else if info.IsDir() {
		fileType = "directory"
	}

	return fileType
}

func visit(files *[]FileInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warnf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		fileName := info.Name()
		modTime := info.ModTime()
		fileType := getType(info)

		if !info.IsDir() {
			fileinfo := FileInfo{fileName, modTime, fileType}

			log.WithFields(log.Fields{
				"fileName": fileName,
				"modTime":  modTime,
			}).Debug("enter this file")
			*files = append(*files, fileinfo)
		} else {
			log.WithFields(log.Fields{
				"fileName": fileName,
				"modTime":  modTime,
			}).Debug("enter this directory")
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

func find(word, root string) []FileInfo {
	files := WalkSearch(root)
	var matchedFiles []FileInfo
	for _, file := range files {
		if strings.Contains(file.FileName, word) {
			matchedFiles = append(matchedFiles, file)
			log.WithFields(log.Fields{
				"fileName": file.FileName,
				"modTime":  file.ModTime,
			}).Info("found this file matched")
		}
	}

	return matchedFiles
}
