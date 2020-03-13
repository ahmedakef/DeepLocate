package osutils

import (
	"os"
	"syscall"
	"time"
)

//GetFileInfo return file metadata sturct using it's path
func GetFileInfo(fileinfo os.FileInfo, path string) FileMetadata {
	stat := fileinfo.Sys().(*syscall.Win32FileAttributeData)
	var file FileMetadata = NewFileMetadata(fileinfo)
	file.Path = path + "\\" + file.Name
	file.ATime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	file.CTime = time.Unix(0, stat.CreationTime.Nanoseconds())
	file.MTime = time.Unix(0, stat.LastWriteTime.Nanoseconds())
	return file
}
