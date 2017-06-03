package main

import (
	"io/ioutil"
	"os"
	"strings"
)

var videoFileExtensions = [...]string{".mp4", ".mov", ".mkv", ".avi", ".flv"}

func isVideoFile(filename string) bool {
	for _, extension := range videoFileExtensions {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}

	return false
}

func dirExists(dirname string) bool {
	st, err := os.Stat(dirname)
	if err != nil {
		return false
	}

	return st.IsDir()
}

func getVideoListFromDir(dirpath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return nil, err
	}

	list := make([]string, 0, 8)
	for _, f := range files {
		if isVideoFile(f.Name()) {
			list = append(list, f.Name())
		}
	}

	return list, nil
}
