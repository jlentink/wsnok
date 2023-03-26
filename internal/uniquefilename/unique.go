package uniquefilename

import (
	"errors"
	"github.com/spf13/afero"
	"os"
	"path"
	"strings"
)

var fs = afero.NewOsFs()

// GetUniqueFilename finds a free filename.
func GetUniqueFilename(filename string, OverWrite bool) string {
	if _, err := fs.Stat(filename); errors.Is(err, os.ErrNotExist) || OverWrite {
		return filename
	}
	return getFilenameIncrement(NewFile(filename))
}

func GetUniqueFilenameFromUrl(url string, OverWrite bool) string {
	if strings.Contains(url, "?") {
		url = url[:strings.Index(url, "?")]
	}
	filename := path.Base(url)
	return GetUniqueFilename(filename, OverWrite)
}

func getFilenameIncrement(file *File) string {
	cFilename := file.NextName()
	if _, err := fs.Stat(cFilename); errors.Is(err, os.ErrNotExist) {
		return cFilename
	}
	return getFilenameIncrement(file)
}
