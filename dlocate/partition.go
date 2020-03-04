package main

import (
	"fmt"
	"os"
)

// Partition conatins basic info about partitions
type Partition struct {
	Index       int
	Root        string
	Directories []string
	FilesNumber int
	Children    []int
	Extenstion  SignatureFile
	ExtenstionH SignatureFile
	//TODO implement versioning
	//TODO add metadata partition pointer
	//TODO add content partition pointer
}

// NewPartition creates new partition with 0 files and 0 children
func NewPartition(index int, root string) Partition {
	return Partition{Index: index, Root: root, FilesNumber: 0, Extenstion: newSignatureFile(), ExtenstionH: newSignatureFileH()}
}

func (p *Partition) addDir(path string) {
	files := ListFiles(path)

	cnt := 0
	p.Directories = append(p.Directories, p.getRelativePath(path))
	for _, file := range files {
		if !file.IsDir() {
			cnt++
			p.addExtenstion(file)
		}
		//TODO update metadata and content partitions
	}
	p.FilesNumber += cnt
	p.ExtenstionH.or(p.Extenstion)
}

func (p *Partition) addExtenstion(file os.FileInfo) {
	extension := getFileExtenstion(file.Name())
	index := getExtenstionIndex(extension)
	p.Extenstion.setBit(index)
}

func (p *Partition) getRelativePath(path string) string {
	return path[len(p.Root):]
}

func (p *Partition) addChild(c *Partition) {
	p.ExtenstionH.or(c.ExtenstionH)
	p.Children = append(p.Children, c.Index)
}

func (p *Partition) hasExtenstion(index int) bool {
	return p.Extenstion.getBit(index)
}

func (p *Partition) hasExtenstionH(index int) bool {
	return p.ExtenstionH.getBit(index)
}

func (p *Partition) printPartition() {
	fmt.Println(p.Index)
	fmt.Println(p.Root)
	fmt.Println(p.FilesNumber)
	fmt.Println(p.Children)
	//fmt.Println(p.Directories)
	//fmt.Println(p.Extenstion)
	//fmt.Println(p.ExtenstionH)
}
