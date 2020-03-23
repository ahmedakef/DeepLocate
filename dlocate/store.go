package main

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	structure "./dataStructures"
	log "github.com/Sirupsen/logrus"
)

func readGob(path string, object interface{}) error {
	dataFile, err := os.Open(filepath.FromSlash(path))

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	// ensure to close the file after the fuction end
	defer dataFile.Close()

	var buf io.Reader = dataFile
	//buf, _ = gzip.NewReader(dataFile)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(object)

	if err != nil {
		return err
	}
	return nil
}

func saveGob(path string, object interface{}) error {

	// FromSlash converts / to the specific file system separator
	dataFile, err := os.Create(filepath.FromSlash(path))
	if err != nil {
		return err
	}
	defer dataFile.Close()

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(object)

	return nil
}

func savePartitionGob(partition *Partition) {
	partition.Root = filepath.ToSlash(partition.Root)
	partitionsPath := filepath.FromSlash("indexFiles/partitions/")
	if _, err := os.Stat(partitionsPath); os.IsNotExist(err) {
		os.MkdirAll(partitionsPath, os.ModePerm)
	}

	path := "indexFiles/partitions/p" + strconv.Itoa(partition.Index) + ".gob"
	err := saveGob(path, partition)

	if err != nil {
		log.Errorf("Error while storing index for partition %v: %v\n", partition.Index, err)
		os.Exit(1)
	}
}

func readPartitionGob(index int) Partition {
	path := "indexFiles/partitions/p" + strconv.Itoa(index) + ".gob"

	var partition Partition
	err := readGob(path, &partition)
	if err != nil {
		log.Errorf("Error while reading index for partition %q: %v\n", index, err)
		os.Exit(1)
	}
	partition.Root = filepath.FromSlash(partition.Root)
	return partition
}

func saveDirectoryPartition(directoryPartition *DirectoryPartition) {
	path := "indexFiles/directoryPartition.gob"
	err := saveGob(path, directoryPartition)

	if err != nil {
		log.Errorf("Error while creating directoryPartitio file")
		os.Exit(1)
	}
}

func readDirectoryPartitionGob() DirectoryPartition {
	path := "indexFiles/directoryPartition.gob"

	var directoryPartition DirectoryPartition
	err := readGob(path, &directoryPartition)
	if err != nil {
		log.Error("Error while reading directoryPartition")
		os.Exit(1)
	}

	return directoryPartition
}

func savePartitionFilesGob(partitionIndex int, partitionFiles []string) {

	path := "indexFiles/filepaths/f" + strconv.Itoa(partitionIndex) + ".gob"
	err := saveGob(path, partitionFiles)

	if err != nil {
		log.Errorf("Error while creating partitionfiles file")
		os.Exit(1)
	}
}

func readPartitionFilesGob(partitionIndex int) []string {
	// TODO return error and results and check the error outside
	path := "indexFiles/filepaths/f" + strconv.Itoa(partitionIndex) + ".gob"

	var partitionFiles []string
	err := readGob(path, &partitionFiles)
	if err != nil {
		log.Error("Error while reading partitionfiles")
		os.Exit(1)
	}

	return partitionFiles
}

func savePartitionMetaGob(partitionIndex int, tree structure.KDTree) {
	path := "indexFiles/metadata/m" + strconv.Itoa(partitionIndex) + ".gob"
	err := saveGob(path, tree)

	if err != nil {
		log.Errorf("Error while creating files metadata tree")
		os.Exit(1)
	}
}

func readPartitionMetaGob(partitionIndex int) structure.KDTree {
	// TODO return error and results and check the error outside
	path := "indexFiles/metadata/m" + strconv.Itoa(partitionIndex) + ".gob"

	var tree structure.KDTree
	err := readGob(path, &tree)
	if err != nil {
		log.Error("Error while reading files metadata tree")
		os.Exit(1)
	}
	return tree
}

// SaveAsJSON save aby datatype as json for better reading while debugging
func SaveAsJSON(data interface{}, filePath string) {
	// create folder if not exits
	lastslash := strings.LastIndex(filePath, "/")
	directoryPath := filePath[:lastslash]
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		os.MkdirAll(directoryPath, os.ModePerm)
	}

	b, _ := json.MarshalIndent(data, "", "\t")
	_ = ioutil.WriteFile(filepath.ToSlash(filePath), b, 0644)
}
