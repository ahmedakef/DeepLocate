package main

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func savePartitionGob(partition *Partition) {
	partition.Root = filepath.ToSlash(partition.Root)
	partitionsPath := "indexFiles/partitions/"
	if _, err := os.Stat(partitionsPath); os.IsNotExist(err) {
		os.MkdirAll(partitionsPath, os.ModePerm)
	}

	// FromSlash converts / to the specific file system separator
	dataFile, err := os.Create(filepath.FromSlash("indexFiles/partitions/p" + strconv.Itoa(partition.Index) + ".gob"))

	if err != nil {
		log.Errorf("Error while storing index for partition %v: %v\n", partition.Index, err)
		os.Exit(1)
	}

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(partition)

	dataFile.Close()
}

func readPartitionGob(index int) Partition {
	dataFile, err := os.Open(filepath.FromSlash("indexFiles/partitions/p" + strconv.Itoa(index) + ".gob"))

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	// ensure to close the file after the fuction end
	defer dataFile.Close()

	var buf io.Reader = dataFile
	//buf, _ = gzip.NewReader(dataFile)
	dec := gob.NewDecoder(buf)

	var partition Partition

	err = dec.Decode(&partition)

	if err != nil {
		log.Errorf("Error while reading index for partition %q: %v\n", index, err)
		os.Exit(1)
	}

	partition.Root = filepath.FromSlash(partition.Root)
	return partition
}

func saveDirectoryPartition(directoryPartition *DirectoryPartition) {

	directoryPartitionPath := "indexFiles/directoryPartition.gob"

	// FromSlash converts / to the specific file system separator
	dataFile, err := os.Create(filepath.FromSlash(directoryPartitionPath))

	if err != nil {
		log.Errorf("Error while creating directoryPartitio file")
		os.Exit(1)
	}
	defer dataFile.Close()

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(directoryPartition)

}

func readDirectoryPartitionGob(index int) DirectoryPartition {
	dataFile, err := os.Open(filepath.FromSlash("indexFiles/partitions/directoryPartition.gob"))

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	// ensure to close the file after the fuction end
	defer dataFile.Close()

	var buf io.Reader = dataFile
	dec := gob.NewDecoder(buf)

	var directoryPartition DirectoryPartition

	err = dec.Decode(&directoryPartition)

	if err != nil {
		log.Error("Error while reading directoryPartition")
		os.Exit(1)
	}

	return directoryPartition
}

// SaveAsJSON save aby datatype as json for better reading while debugging
func SaveAsJSON(Data interface{}, FilePath string) {
	b, _ := json.MarshalIndent(Data, "", "\t")
	_ = ioutil.WriteFile(filepath.ToSlash(FilePath), b, 0644)
}
