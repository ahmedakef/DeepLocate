package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

var tokens = make(chan struct{}, 20)
var wg sync.WaitGroup

// getPartitionFiles gets files of partition and its children
func getPartitionsFiles(partitions []int, path string, fileNamesChan chan string) {

	for _, partitionIndex := range partitions {
		partition := indexInfo.getPartition(partitionIndex)

		// check that partition is related to the path
		if !partition.inSameDirection(path) {
			// return []string{}
			return
		}

		var partitionFiles map[string][]string
		wg.Add(1)
		go func() {
			defer wg.Done()
			tokens <- struct{}{} // acuire token
			partitionFiles = partition.getPartitionFiles()
			<-tokens // release the token

			// fileNames := make([]string, partition.FilesNumber+len(partitionFiles))

			for path, files := range partitionFiles {
				for _, fileName := range files {
					fileNamesChan <- partition.Root + path + fileName
					// fmt.Println(partitionIndex)
					// fmt.Println(partition.Root + path + fileName)
					// fileNames[i] = partition.Root + path + fileName
					// i++
				}
				// search in directoris also
				fileNamesChan <- partition.Root + path[:len(path)-1]
				// fileNames[i] = partition.Root + path[:len(path)-1]
				// i++
			}
			// os.Exit(1)
		}()
	}

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

func findInFileNames(query string, fileNamesChan chan string, scores chan pair) {
	if query == "" {
		for fileName := range fileNamesChan {
			scores <- pair{fileName, 1}
		}
	}

	words := strings.Fields(query)

	for fileName := range fileNamesChan {
		lastslash := strings.LastIndex(fileName, "/")
		exactFileName := fileName[lastslash+1:]
		score := pair{fileName, 0}

		// exact match has high score
		if exactFileName == query {
			score.second += 10
		}

		// partial match
		for _, word := range words {
			if strings.Contains(exactFileName, word) {
				score.second++
			}
		}
		if score.second > 0 {
			scores <- score
		}

	}

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

	children := getPartitionClildren(partitionIndex, path)
	fileNamesChan := make(chan string)

	getPartitionsFiles(children, path, fileNamesChan)

	// close(fileNamesChan)
	var matchedFiles = []string{}

	scores := make(chan pair)

	go findInFileNames(query, fileNamesChan, scores)

	go func() {
		wg.Wait()
		close(fileNamesChan)
		close(scores)

	}()

	// score maybe used later to sort of filter first n entries
	for score := range scores {
		log.WithFields(log.Fields{
			"fileName": score.first,
		}).Debug("found this file matched")
		fmt.Println(score.first)
		matchedFiles = append(matchedFiles, score.first)

	}

	var contentMatchedFiles = []string{}

	if searchContent {
		log.Debug("Start searching file content ...")
		fmt.Println("=====\nStart searching file content ...\n=====")

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

	var fileNames []string
	fileNamesChan := make(chan string)
	log.Info("Start searching file metadata ...")

	//get parition index:
	for _, partitionIndex := range partitions {
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
				fileNamesChan <- file.Path
				continue
			}
			for _, extention := range extentions {
				if extention == file.Extension {
					fileNames = append(fileNames, file.Path)
					fileNamesChan <- file.Path
					break
				}
			}
		}
	}

	var matchedFiles = []string{}

	scores := make(chan pair)
	findInFileNames(query, fileNamesChan, scores)
	// score maybe used later to sort of filter first n entries
	for score := range scores {
		log.WithFields(log.Fields{
			"fileName": score.first,
		}).Info("found this file matched")
		matchedFiles = append(matchedFiles, score.first)

	}

	var contentMatchedFiles = []string{}

	if searchContent {
		log.Info("Start searching file content ...")

		contentMatchedFiles = invertedIndex.SearchIn(partitions, query, -1, fileNames)

		for _, fileName := range contentMatchedFiles {
			log.WithFields(log.Fields{
				"fileName": fileName,
			}).Info("found this file content matched")

		}

	}

	return matchedFiles, contentMatchedFiles
}

type pair struct {
	first  string
	second int
}
