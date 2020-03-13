package main

import log "github.com/Sirupsen/logrus"

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
		if !file.IsDir {
			cnt++
			p.addExtension(file.Extension)
		}

		//TODO update metadata and content partitions
	}
	p.FilesNumber += cnt
	p.ExtenstionH.or(p.Extenstion)
}

func (p *Partition) addExtension(extension string) {
	index := getExtensionIndex(extension)
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
