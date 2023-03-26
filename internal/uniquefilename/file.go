package uniquefilename

import (
	"path/filepath"
	"strconv"
)

type File struct {
	Increment int
	Filename  string
	Name      string
	Extension string
}

func (f *File) NextName() string {
	f.Increment++
	return f.Name + strconv.Itoa(f.Increment) + "." + f.Extension
}

func NewFile(fileName string) *File {
	return &File{
		Filename:  fileName,
		Name:      fileName[:len(fileName)-len(filepath.Ext(fileName))],
		Extension: filepath.Ext(fileName)[1:],
		Increment: 0,
	}
}
