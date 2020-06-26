package structures

import (
	"sort"
	"strconv"

	utils "dlocate/osutils"
	python "dlocate/python"

	log "github.com/sirupsen/logrus"
)

//InvertedIndex is a 2 level inverted index to search file contents
type InvertedIndex struct {
	content             map[string]map[int]map[int]float32 //word, parition, file index, freq(score)
	filesToIndices      map[int]map[int]string             //parition, file index, filepath
	NextFileID          map[int]int                        // partiion filepath
	contentCache        Cache
	filesToIndicesCache Cache
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
		invertedIndex.loadPartitionInvertedIndex(partition, word)
		invertedIndex.content[word][partition][fileID] = freq
		invertedIndex.savePartitionInvertedIndex(partition, word)
	}
}

//Search the content index for a title (put limit = -1 for all results)
func (invertedIndex *InvertedIndex) Search(partitions []int, query string, limit int) []string {

	var words []string
	python.ExecuteScript("keyword_extraction/text_cleaning.py", query, &words)

	scores := make(map[pair]float32)

	for _, word := range words {

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

	if len(values) == 0 {
		return []string{}
	}

	limitIndex := len(values) - limit
	if limitIndex < 0 {
		limitIndex = 0
	}
	scoreLimit := values[limitIndex]

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

//SearchIn searchs the content index for a query but within a specific files only
func (invertedIndex *InvertedIndex) SearchIn(partitions []int, query string, limit int, within []string) []string {
	var words []string
	python.ExecuteScript("keyword_extraction/text_cleaning.py", query, &words)

	scores := make(map[string]float32)
	for _, file := range within {
		scores[file] = 0
	}

	for _, partition := range partitions {
		invertedIndex.loadPartitionDir(partition)
		for _, word := range words {
			invertedIndex.loadPartitionInvertedIndex(partition, word)
			for fileID, score := range invertedIndex.content[word][partition] {
				_, ok := scores[invertedIndex.filesToIndices[partition][fileID]]
				if ok {
					scores[invertedIndex.filesToIndices[partition][fileID]] += score
				}
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

	if len(values) == 0 {
		return []string{}
	}

	limitIndex := len(values) - limit
	if limitIndex < 0 {
		limitIndex = 0
	}
	scoreLimit := values[limitIndex]

	results := []string{}
	for file, score := range scores {
		if score >= scoreLimit {
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
	invertedIndex.contentCache = GetCache(100)
	invertedIndex.filesToIndicesCache = GetCache(10)
}

func (invertedIndex *InvertedIndex) loadPartitionDir(partitionID int) {
	pDir, ok := invertedIndex.filesToIndicesCache.Get(strconv.Itoa(partitionID))
	if ok {
		invertedIndex.filesToIndices[partitionID] = pDir.(map[int]string)
		return
	}

	path := "indexFiles/content/filesToIndices" + strconv.Itoa(partitionID) + ".gob"

	var partitionDir map[int]string
	err := utils.ReadGob(path, &partitionDir)
	if err != nil {
		invertedIndex.filesToIndices[partitionID] = make(map[int]string)
		return
	}
	invertedIndex.filesToIndices[partitionID] = partitionDir
}

func (invertedIndex *InvertedIndex) savePartitionDir(partitionID int) {
	path := "indexFiles/content/filesToIndices" + strconv.Itoa(partitionID) + ".gob"
	err := utils.SaveGob(invertedIndex.filesToIndices[partitionID], path)
	if err != nil {
		log.Error(err)
	}
	invertedIndex.filesToIndicesCache.Set(strconv.Itoa(partitionID), invertedIndex.filesToIndices[partitionID])
}

func (invertedIndex *InvertedIndex) loadPartitionInvertedIndex(partition int, keyword string) {
	key := keyword + "-" + strconv.Itoa(partition)

	if invertedIndex.content[keyword] == nil {
		invertedIndex.content[keyword] = make(map[int]map[int]float32)
	}

	pInvertedIndex, ok := invertedIndex.contentCache.Get(key)

	if ok {
		invertedIndex.content[keyword][partition] = pInvertedIndex.(map[int]float32)
		return
	}

	path := "indexFiles/content/" + keyword + "-" + strconv.Itoa(partition) + ".gob"
	var parInvertedIndex map[int]float32

	err := utils.ReadGob(path, &parInvertedIndex)

	if err != nil {
		invertedIndex.content[keyword][partition] = make(map[int]float32)
		return
	}

	invertedIndex.content[keyword][partition] = parInvertedIndex
}

func (invertedIndex *InvertedIndex) savePartitionInvertedIndex(partition int, keyword string) {
	path := "indexFiles/content/" + keyword + "-" + strconv.Itoa(partition) + ".gob"

	err := utils.SaveGob(invertedIndex.content[keyword][partition], path)
	if err != nil {
		log.Error(err)
	}
	invertedIndex.contentCache.Delete(keyword + "-" + strconv.Itoa(partition))
}

type pair struct {
	first, second int
}
