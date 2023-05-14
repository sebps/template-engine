package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func ReadFileContent(filePath string) (string, error) {
	fileRawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(fileRawContent), nil
}

func WriteFileContent(path string, content string) error {
	// prepare dir hierarchy
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	// write content
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)

	return err
}
