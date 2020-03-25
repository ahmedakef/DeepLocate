package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	structure "./dataStructures"
	utils "./osutils"
	log "github.com/Sirupsen/logrus"
)

// Partition conatins basic info about partitions
type Partition struct {
	Index        int
	Root         string
	Directories  []string
	FilesNumber  int
	Children     []int // index for children partitions
	Extenstion   SignatureFile
	ExtenstionH  SignatureFile
	filePaths    []string
	metadataTree structure.KDTree
	//TODO implement versioning
	//TODO add metadata partition pointer
	//TODO add content partition pointer
}

// NewPartition creates new partition with 0 files and 0 children
func NewPartition(index int, root string) Partition {
	return Partition{Index: index, Root: root, FilesNumber: 0, Extenstion: newSignatureFile(), ExtenstionH: newSignatureFileH()}
}

func (p *Partition) addDir(path string) {
	files := utils.ListFiles(path)
	cnt := 0
	p.Directories = append(p.Directories, p.getRelativePath(path))
	for _, file := range files {
		if !file.IsDir {
			p.filePaths = append(p.filePaths, p.getRelativePath(file.Path))
			p.metadataTree.Insert(&file)
			cnt++
			p.addExtension(file.Extension)
		}
		//TODO content partitions
	}
	p.FilesNumber += cnt
	p.ExtenstionH.or(p.Extenstion)
}

func (p *Partition) addExtension(extension string) {
	index := indexInfo.getExtensionIndex(extension)
	p.Extenstion.setBit(index)
}

// dirInDirs ensures that root exits inside any of directories
func (p *Partition) containsDir(root string) bool {
	for _, dir := range p.Directories {
		if strings.HasPrefix(dir+"/", root+"/") {
			return true
		}
	}
	return false
}

// inSameDirection ensures that path is parent or child of the partition path
func (p *Partition) inSameDirection(path string) bool {
	if len(path) > len(p.Root) {
		return strings.HasPrefix(path, p.Root)
	}
	return strings.HasPrefix(p.Root, path)
}

func (p *Partition) getRelativePath(path string) string {
	// if root is longer then it is one of path's children
	if len(p.Root) > len(path) {
		return ""
	}
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

func (p *Partition) saveAsGob() {
	p.Root = filepath.ToSlash(p.Root)
	partitionsPath := filepath.FromSlash("indexFiles/partitions/")
	if _, err := os.Stat(partitionsPath); os.IsNotExist(err) {
		os.MkdirAll(partitionsPath, os.ModePerm)
	}

	path := "indexFiles/partitions/p" + strconv.Itoa(p.Index) + ".gob"
	err := saveGob(path, p)

	if err != nil {
		log.Errorf("Error while storing index for partition %v: %v\n", p.Index, err)
		os.Exit(1)
	}

}

func (p *Partition) printPartition() {
	log.WithFields(log.Fields{
		"Index":       p.Index,
		"Root":        p.Root,
		"FilesNumber": p.FilesNumber,
		"Children":    p.Children,
		// "Directories": p.Directories,
		// "Extenstion":  p.Extenstion,
		// "ExtenstionH": p.ExtenstionH,
	}).Infof("partition %v info :", p.Index)

}
