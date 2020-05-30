package structures

import (
	"sort"
	"strconv"
	"strings"

	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

//InvertedIndex is a 2 level inverted index to search file contents
type InvertedIndex struct {
	content        map[string]map[int]map[int]float32 //word, parition, file index, freq(score)
	filesToIndices map[int]map[int]string             //parition, file index, filepath
	NextFileID     map[int]int                        // partiion filepath
	//TODO replace filepath with file relative path
}

//Insert add[update] a file contents to[in] the content index
func (invertedIndex *InvertedIndex) Insert(partition int, file string, contents map[string]float32) {
	invertedIndex.loadPartitionDir(partition)
	fileID := -1
	for k, v := range invertedIndex.filesToIndices[partition] {
		if v == file {
			fileID = k
		}
	}
	if fileID == -1 {
		fileID = invertedIndex.NextFileID[partition]
		invertedIndex.NextFileID[partition]++
		invertedIndex.filesToIndices[partition][fileID] = file
		invertedIndex.savePartitionDir(partition)
	}
	for word, freq := range contents {
		if isFiltered(word) {
			continue
		}
		word = stem(word)
		invertedIndex.loadPartitionInvertedIndex(partition, word)
		invertedIndex.content[word][partition][fileID] = freq
		invertedIndex.savePartitionInvertedIndex(partition, word)
	}
}

//Search the content index for a title (put limit = -1 for all results)
func (invertedIndex *InvertedIndex) Search(partitions []int, query string, limit int) []string {
	words := strings.SplitN(query, " ", 1)
	scores := make(map[pair]float32)

	for _, word := range words {
		if isFiltered(word) {
			continue
		}
		word = stem(word)
		for _, partition := range partitions {
			invertedIndex.loadPartitionInvertedIndex(partition, word)
			files := invertedIndex.content[word][partition]
			for fileID, score := range files {
				scores[pair{partition, fileID}] += score
			}
		}
	}

	values := []float32{}
	for _, val := range scores {
		values = append(values, val)
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	if limit == -1 {
		limit = len(values)
	}

	scoreLimit := values[len(values)-limit]

	results := []string{}
	for fileKey, score := range scores {
		if score >= scoreLimit {
			invertedIndex.loadPartitionDir(fileKey.first)
			file := invertedIndex.filesToIndices[fileKey.first][fileKey.second]
			results = append(results, file)
			if len(results) >= limit {
				break
			}
		}
	}
	return results
}

//Save the inverted Index object
func (invertedIndex *InvertedIndex) Save() {
	path := "indexFiles/invertedIndex.gob"
	err := utils.SaveGob(&invertedIndex, path)
	if err != nil {
		log.Error(err)
	}
}

//Load the inverted Index object
func (invertedIndex *InvertedIndex) Load() {
	path := "indexFiles/invertedIndex.gob"
	invertedIndex.content = make(map[string]map[int]map[int]float32)
	invertedIndex.filesToIndices = make(map[int]map[int]string)
	err := utils.ReadGob(path, &invertedIndex)
	if err != nil {
		invertedIndex.NextFileID = make(map[int]int)
	}
}

func (invertedIndex *InvertedIndex) loadPartitionDir(partitionID int) {
	if invertedIndex.filesToIndices[partitionID] != nil {
		return
	}
	path := "indexFiles/content/filesToIndices" + strconv.Itoa(partitionID) + ".gob"
	var partitionDir map[int]string
	err := utils.ReadGob(path, &partitionDir)
	if err != nil {
		invertedIndex.filesToIndices[partitionID] = partitionDir
	}
	invertedIndex.filesToIndices[partitionID] = make(map[int]string)
}

func (invertedIndex *InvertedIndex) savePartitionDir(partitionID int) {
	path := "indexFiles/content/filesToIndices" + strconv.Itoa(partitionID) + ".gob"
	err := utils.SaveGob(invertedIndex.filesToIndices[partitionID], path)
	if err != nil {
		log.Error(err)
	}
}

func (invertedIndex *InvertedIndex) loadPartitionInvertedIndex(partition int, keyword string) {
	if invertedIndex.content[keyword][partition] != nil {
		return
	}
	path := "indexFiles/content/" + keyword + "-" + strconv.Itoa(partition) + ".gob"

	var pInvertedIndex map[int]float32
	err := utils.ReadGob(path, &pInvertedIndex)

	if invertedIndex.content[keyword] == nil {
		invertedIndex.content[keyword] = make(map[int]map[int]float32)
	}

	if err != nil {
		invertedIndex.content[keyword][partition] = pInvertedIndex
	}

	invertedIndex.content[keyword][partition] = make(map[int]float32)
}

func (invertedIndex *InvertedIndex) savePartitionInvertedIndex(partition int, keyword string) {
	path := "indexFiles/content/" + keyword + "-" + strconv.Itoa(partition) + ".gob"

	log.Debug(path)

	err := utils.SaveGob(invertedIndex.content[keyword][partition], path)
	if err != nil {
		log.Error(err)
	}
}

//TODO stem words
func stem(word string) string {
	newWord := word
	return newWord
}

//TODO filter words
func isFiltered(word string) bool {
	return false
}

type pair struct {
	first, second int
}
