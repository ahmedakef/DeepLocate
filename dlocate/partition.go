package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Partition conatins basic info about partitions
type Partition struct {
	index       int
	root        string
	directories []string
	filesNumber int
	children    []*Partition
	extenstion  SignatureFile
	extenstionH SignatureFile
	//TODO implement versioning
	//TODO add metadata partition pointer
	//TODO add content partition pointer
}

// NewPartition creates new partition with 0 files and 0 children
func NewPartition(index int, root string) Partition {
	return Partition{index: index, root: root, filesNumber: 0, extenstion: newSignatureFile(), extenstionH: newSignatureFileH()}
}

func (p Partition) addDir(relativePath string) {
	path := p.root + "/" + relativePath

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	p.filesNumber += len(files)
	p.directories = append(p.directories, relativePath)
	for _, file := range files {
		fmt.Println(file.Name())
		//TODO add extenstions to signature file and hierarchical signature files
		//TODO update metadata and content partitions
	}
}

func (p Partition) addChild(c *Partition) {
	p.extenstionH.or(c.extenstionH)
	p.children = append(p.children, c)
}

func (p Partition) hasExtenstion(index int) bool {
	return p.extenstion.getBit(index)
}

func (p Partition) hasExtenstionH(index int) bool {
	return p.extenstionH.getBit(index)
}
