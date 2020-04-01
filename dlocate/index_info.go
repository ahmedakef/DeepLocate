package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

//IndexInfo stores general data about the whole index
type IndexInfo struct {
	CurIndex          int
	CurExtensionIndex int
	Extensions        map[string]int
	Roots             []int
	partitions        map[int]*Partition
}

func getIndexInfo() IndexInfo {
	path := "indexFiles/indexInfo.gob"

	var indexInfo IndexInfo
	err := readGob(path, &indexInfo)
	if err != nil {
		return IndexInfo{CurIndex: 0, CurExtensionIndex: -1,
			Extensions: make(map[string]int), partitions: make(map[int]*Partition),
			Roots: make([]int, 0)}
	}
	indexInfo.partitions = make(map[int]*Partition)
	return indexInfo
}

func (x *IndexInfo) loadRoots() {
	for _, root := range x.Roots {
		partition := readPartitionGob(root)
		x.partitions[root] = &partition
	}
}

func (x *IndexInfo) getExtensionIndex(name string) int {
	index, ok := x.Extensions[name]
	if !ok {
		index = x.getNextExtensionIndex()
		x.Extensions[name] = index
	}
	return index
}

func (x *IndexInfo) getNextExtensionIndex() int {
	x.CurExtensionIndex++
	return x.CurExtensionIndex
}

func (x *IndexInfo) getNextPartitionIndex() int {
	x.CurIndex++
	return x.CurIndex
}

func (x *IndexInfo) addRoot(path string) *Partition {
	root := NewPartition(x.getNextPartitionIndex(), path)
	x.Roots = append(x.Roots, root.Index)
	x.addPartition(&root)
	return &root
}

func (x *IndexInfo) addPartition(p *Partition) {
	x.partitions[p.Index] = p
}

func (x *IndexInfo) getPartition(index int) *Partition {
	partition, ok := x.partitions[index]
	if !ok {
		p := readPartitionGob(index)
		x.addPartition(&p)
		return &p
	}
	return partition
}

func (x *IndexInfo) saveAsGob() {
	SaveAsJSON(x, "indexFiles/indexInfo.json")
	path := "indexFiles/indexInfo.gob"
	err := saveGob(x, path)

	if err != nil {
		log.Errorf("Error while creating indexInfo file")
		os.Exit(1)
	}
}

func (x *IndexInfo) savePartitions() {
	for _, partition := range x.partitions {
		log.Debugf("start saving Partition %v\n", partition.Index)
		partition.printPartition()

		partition.saveAsGob()

		// save metadata tree inside the partition
		savePartitionMetaGob(partition.Index, partition.metadataTree)
	}
}

func (x *IndexInfo) removeRoot(root int) {
	for i, value := range x.Roots {
		if value == root {
			x.Roots[i] = x.Roots[len(x.Roots)-1]
			x.Roots = x.Roots[:len(x.Roots)-1]
			break
		}
	}
}
