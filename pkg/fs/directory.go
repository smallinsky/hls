package fs

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var videoExtensions = []string{
	".mp4",
	".mov",
	".mkv",
	".avi",
	".flv",
}

func IsVideoFile(filename string) bool {
	for _, ext := range videoExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func DirExist(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func Mk(path string) error {
	const prem = 0755
	return os.MkdirAll(path, prem)
}

func ListVideoFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, 8)
	for _, f := range files {
		if IsVideoFile(f.Name()) {
			list = append(list, f.Name())
		}
	}
	return list, nil
}

func FileName(file string) string {
	return strings.TrimSuffix(path.Base(file), path.Ext(file))
}

func FileExt(file string) string {
	return path.Ext(file)
}
