package main

import (
	"fmt"
	"os"
)

// Partition conatins basic info about partitions
type Partition struct {
	index       int
	root        string
	directories []string
	filesNumber int
	children    []int
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

func (p *Partition) addDir(path string) {
	files := ListFiles(path)

	cnt := 0
	p.directories = append(p.directories, p.getRelativePath(path))
	for _, file := range files {
		if !file.IsDir() {
			cnt++
			p.addExtenstion(file)
		}
		//TODO update metadata and content partitions
	}
	p.filesNumber += cnt
	p.extenstionH.or(p.extenstion)
}

func (p *Partition) addExtenstion(file os.FileInfo) {
	extension := getFileExtenstion(file.Name())
	index := getExtenstionIndex(extension)
	p.extenstion.setBit(index)
}

func (p *Partition) getRelativePath(path string) string {
	return path[len(p.root):]
}

func (p *Partition) addChild(c *Partition) {
	p.extenstionH.or(c.extenstionH)
	p.children = append(p.children, c.index)
}

func (p *Partition) hasExtenstion(index int) bool {
	return p.extenstion.getBit(index)
}

func (p *Partition) hasExtenstionH(index int) bool {
	return p.extenstionH.getBit(index)
}

func (p *Partition) printPartition() {
	fmt.Println(p.index)
	fmt.Println(p.root)
	fmt.Println(p.filesNumber)
	fmt.Println(p.children)
	//fmt.Println(p.directories)
	//fmt.Println(p.extenstion)
	//fmt.Println(p.extenstionH)
}
