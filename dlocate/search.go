package main

import (
	"fmt"
	"strconv"
	"strings"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

// getPartitionFiles gets files of partition and its children
func getPartitionFiles(partitionIndex int, path string) []string {
	partition := indexInfo.getPartition(partitionIndex)

	// check that partition is related to the path
	if !partition.inSameDirection(path) {
		return []string{}
	}

	partitionFiles := partition.getPartitionFiles()

	fileNames := make([]string, partition.FilesNumber+len(partitionFiles))
	i := 0
	for directory, files := range partitionFiles {
		for _, fileName := range files {
			if strings.HasPrefix(partition.Root+directory, path+"/") {
				fileNames[i] = partition.Root + directory + fileName
				i++
			}
		}
		// search in directoris also
		if strings.HasPrefix(partition.Root+directory, path+"/") {
			fileNames[i] = partition.Root + directory[:len(directory)-1]
			i++
		}
	}

	for _, child := range partition.Children {
		// go func() {
		// 	fileNames = append(fileNames, getPartitionFiles(child, path)...)
		// }()
		fileNames = append(fileNames, getPartitionFiles(child, path)...)
	}

	return fileNames
}

// getPartitionClildren get the children of partioin that is related
// to the given path (either parent or child)
// and excludes the non relevant partitions
func getPartitionClildren(partitionIndex int, path string) []int {
	partition := indexInfo.getPartition(partitionIndex)

	// check that partition is related to the path
	if !partition.inSameDirection(path) {
		return []int{}
	}

	children := []int{partitionIndex}

	for _, child := range partition.Children {
		children = append(children, getPartitionClildren(child, path)...)
	}

	return children
}

func findInFileNames(query string, fileNames []string) map[string]int {
	scores := make(map[string]int) // map from file path to its score
	if query == "" {
		for _, fileName := range fileNames {
			scores[fileName] = 1
		}
		return scores
	}

	// exact match has high score
	for _, fileName := range fileNames {
		lastslash := strings.LastIndex(fileName, "/")
		exactFileName := fileName[lastslash+1:]
		if exactFileName == query {
			scores[fileName] += 10
		}
	}

	// partial match
	words := strings.Fields(query)
	for _, word := range words {

		for _, fileName := range fileNames {
			lastslash := strings.LastIndex(fileName, "/")
			exactFileName := fileName[lastslash+1:]
			if strings.Contains(exactFileName, word) {
				scores[fileName]++
			}
		}
	}

	return scores
}

// query: word to search
// path: directoy to search in
// searchContent : bool to indicate search content or not
func find(query, path string, searchContent bool) ([]string, []string) {

	log.Debug("Start searching file names ...")

	// get all files names in this path and its children
	partitionIndex := directoryPartition.getPathPartition(path)

	if partitionIndex == -1 {
		return []string{}, []string{}
	}

	fileNames := getPartitionFiles(partitionIndex, path)

	var matchedFiles = []string{}

	scores := findInFileNames(query, fileNames)

	// score maybe used later to sort of filter first n entries
	for fileName := range scores {
		log.WithFields(log.Fields{
			"fileName": fileName,
		}).Debug("found this file matched")
		fmt.Println(fileName)
		matchedFiles = append(matchedFiles, fileName)

	}

	var contentMatchedFiles = []string{}

	if searchContent {
		log.Debug("Start searching file content ...")
		fmt.Println("=====\nStart searching file content ...\n=====")
		children := getPartitionClildren(partitionIndex, path)
		contentMatchedFiles = invertedIndex.Search(children, query, -1)

		for _, fileName := range contentMatchedFiles {
			log.WithFields(log.Fields{
				"fileName": fileName,
			}).Debug("found this file content matched")
			fmt.Println(fileName)

		}

	}
	return matchedFiles, contentMatchedFiles
}

func metaSearch(query, path string, searchContent bool, start utils.FileMetadata,
	end utils.FileMetadata, extentions []string) ([]string, []string) {

	//Size - file size in bytes
	//CTime - change time (last file name or path change)
	//MTime - modify time Max(last content change, CTime)
	//ATime - access time Max(last opened, MTime)

	//time.Date(2019, 1, 1, 20, 34, 58, 651387237, time.UTC)

	// get all partitions in given path
	partitionIndex := directoryPartition.getPathPartition(path)

	if partitionIndex == -1 {
		return []string{}, []string{}
	}

	partitions := getPartitionClildren(partitionIndex, path)

	validPartitions := []int{}

	if len(extentions) == 0 {
		validPartitions = partitions
	} else {
		for _, partitionIndex := range partitions {
			partition := indexInfo.getPartition(partitionIndex)
			for _, extention := range extentions {
				if partition.hasExtension(extention) {
					validPartitions = append(validPartitions, partitionIndex)
					break
				}
			}
		}
	}

	var fileNames []string
	log.Info("Start searching file metadata ...")

	//get parition index:
	for _, partitionIndex := range validPartitions {
		val, ok := indexInfo.metaCache.Get(strconv.Itoa(partitionIndex))
		var tree structure.KDTree
		if !ok {
			tree = readPartitionMetaGob(partitionIndex)
			indexInfo.metaCache.Set(strconv.Itoa(partitionIndex), tree)
		} else {
			tree = val.(structure.KDTree)
		}

		filesInfo := tree.SearchPartial(&start, &end)

		for _, file := range filesInfo {
			// if user hasn't chosed any extention
			if len(extentions) == 0 {
				fileNames = append(fileNames, file.Path)
				continue
			}
			for _, extention := range extentions {
				if extention == file.Extension {
					fileNames = append(fileNames, file.Path)
					break
				}
			}
		}
	}

	var matchedFiles = []string{}

	scores := findInFileNames(query, fileNames)
	// score maybe used later to sort of filter first n entries
	for fileName := range scores {
		log.WithFields(log.Fields{
			"fileName": fileName,
		}).Info("found this file matched")
		matchedFiles = append(matchedFiles, fileName)

	}

	var contentMatchedFiles = []string{}

	if searchContent {
		log.Info("Start searching file content ...")

		contentMatchedFiles = invertedIndex.SearchIn(validPartitions, query, -1, fileNames)

		for _, fileName := range contentMatchedFiles {
			log.WithFields(log.Fields{
				"fileName": fileName,
			}).Info("found this file content matched")

		}

	}

	return matchedFiles, contentMatchedFiles
}
