package main

import (
	log "github.com/Sirupsen/logrus"
)

// word: word to search
// root: directoy to search in
func find(word, root string) []FileInfo {
	var directoryPartition DirectoryPartition

	directoryPartition = readDirectoryPartitionGob()
	partitionIndex := directoryPartition.getDirectoryPartition(root)
	log.Info(partitionIndex)
	return nil
	// TODO read the files in this partition and its children
	// content, err := ioutil.ReadFile("files.json")
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }
	// var files []FileInfo
	// json.Unmarshal(content, &files)

	// var matchedFiles []FileInfo
	// for _, file := range files {
	// 	if strings.Contains(file.FileName, word) &&
	// 		strings.HasPrefix(file.LinuxPath, root) {
	// 		matchedFiles = append(matchedFiles, file)
	// 		log.WithFields(log.Fields{
	// 			"fileName": file.FileName,
	// 			"modTime":  file.ModTime,
	// 		}).Info("found this file matched")
	// 	}
	// }

	// return matchedFiles
}
