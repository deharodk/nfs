package main

import (
	l4g "code.google.com/p/log4go"
	"os"
	"path/filepath"
)

const (
	FILE_TYPE = "*.xml"
)

func ListFiles(globPattern string) (matches []string, err error) {
	return filepath.Glob(globPattern)
}

func GetGlobPatternList(options map[string]interface{}) (output []string) {
	baseDir := options["--path"].(string)
	rfcList, _ := getRFCList(baseDir)

	for _, dir := range rfcList {
		output = append(output, filepath.Join(dir, FILE_TYPE))
		l4g.Debug(filepath.Join(dir, FILE_TYPE))
	}
	return
}

func getRFCList(baseDir string) (matches []string, err error) {
	return filepath.Glob(filepath.Join(baseDir))
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
