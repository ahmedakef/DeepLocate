package osutils

import (
	"os"
	"syscall"
	"time"
)

//GetFileInfo return file metadata sturct using it's path
func GetFileInfo(fileinfo os.FileInfo, path string) FileMetadata {
	var file = NewFileMetadata(fileinfo)
	file.Path = path + "/" + file.Name
	aTime := fileinfo.Sys().(*syscall.Stat_t).Atim
	cTime := fileinfo.Sys().(*syscall.Stat_t).Ctim
	mTime := fileinfo.Sys().(*syscall.Stat_t).Mtim

	file.ATime = time.Unix(aTime.Sec, aTime.Nsec)
	file.CTime = time.Unix(cTime.Sec, cTime.Nsec)
	file.MTime = time.Unix(mTime.Sec, mTime.Nsec)

	return file
}
