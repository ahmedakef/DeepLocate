package osutils

import (
	"os"
	"strings"
	"time"
)

// FileMetadata is a os independant class that describes file metadata
type FileMetadata struct {
	Name      string
	Path      string
	MTime     time.Time
	CTime     time.Time
	ATime     time.Time
	IsDir     bool
	Size      int64
	Extension string
}

// NewFileMetadata creates new file metadata object with basic info from FileInfo class
func NewFileMetadata(fileinfo os.FileInfo) FileMetadata {
	return FileMetadata{Name: fileinfo.Name(), Size: fileinfo.Size(), IsDir: fileinfo.IsDir(), Extension: getFileExtension(fileinfo.Name())}
}

func getFileExtension(name string) string {
	dot := strings.LastIndex(name, ".")
	extention := name[dot+1:]
	return extention
}
