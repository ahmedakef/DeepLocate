package main

import (
	"os"
	"strconv"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

//IndexInfo stores general data about the whole index
type IndexInfo struct {
	CurIndex          int
	CurExtensionIndex int
	Extensions        map[string]int
	Roots             []int
	partitionsCache   structure.Cache
	metaCache         structure.Cache
	filesCache        structure.Cache

	updatedPartitions map[int]*Partition
}

func getIndexInfo() IndexInfo {
	path := "indexFiles/indexInfo.gob"

	var indexInfo IndexInfo
	err := utils.ReadGob(path, &indexInfo)
	if err != nil {
		return IndexInfo{
			CurIndex:          0,
			CurExtensionIndex: -1,
			Extensions:        make(map[string]int),
			partitionsCache:   structure.GetCache(100),
			metaCache:         structure.GetCache(100),
			filesCache:        structure.GetCache(100),
			Roots:             make([]int, 0),
		}
	}
	indexInfo.partitionsCache = structure.GetCache(100)
	indexInfo.metaCache = structure.GetCache(100)
	indexInfo.filesCache = structure.GetCache(100)

	indexInfo.updatedPartitions = make(map[int]*Partition)

	return indexInfo
}

func (indexInfo *IndexInfo) getExtensionIndex(name string) int {
	index, ok := indexInfo.Extensions[name]
	if !ok {
		index = indexInfo.getNextExtensionIndex()
		indexInfo.Extensions[name] = index
	}
	return index
}

func (indexInfo *IndexInfo) getNextExtensionIndex() int {
	indexInfo.CurExtensionIndex++
	return indexInfo.CurExtensionIndex
}

func (indexInfo *IndexInfo) getNextPartitionIndex() int {
	indexInfo.CurIndex++
	return indexInfo.CurIndex
}

func (indexInfo *IndexInfo) addRoot(path string) Partition {
	root := NewPartition(indexInfo.getNextPartitionIndex(), path)
	indexInfo.Roots = append(indexInfo.Roots, root.Index)
	indexInfo.addPartition(root)
	return root
}

func (indexInfo *IndexInfo) addPartition(p Partition) {
	indexInfo.partitionsCache.Set(strconv.Itoa(p.Index), p)
}

func (indexInfo *IndexInfo) getPartition(index int) Partition {
	partition, ok := indexInfo.partitionsCache.Get(strconv.Itoa(index))
	if !ok {
		p := readPartitionGob(index)
		indexInfo.addPartition(p)
		return p
	}
	return partition.(Partition)
}

func (indexInfo *IndexInfo) saveAsGob() {
	path := "indexFiles/indexInfo.gob"
	err := utils.SaveGob(indexInfo, path)

	if err != nil {
		log.Errorf("Error while creating indexInfo file")
		os.Exit(1)
	}
}

func (indexInfo *IndexInfo) clearPartitions() {
	indexInfo.partitionsCache.Clear()
	indexInfo.metaCache.Clear()
	indexInfo.filesCache.Clear()
}

func isRoot(path string) (int, bool) {
	partitionIndex := directoryPartition.getPathPartition(path)

	for _, root := range indexInfo.Roots {
		if partitionIndex == root {
			parition := indexInfo.getPartition(root)
			if parition.Root == path {
				return root, true
			}
		}
	}
	return 0, false
}

func (indexInfo *IndexInfo) removeRoot(root int) {
	for i, value := range indexInfo.Roots {
		if value == root {
			indexInfo.Roots[i] = indexInfo.Roots[len(indexInfo.Roots)-1]
			indexInfo.Roots = indexInfo.Roots[:len(indexInfo.Roots)-1]
			break
		}
	}
}
