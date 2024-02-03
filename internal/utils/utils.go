package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func RegexpEscapeString(str string) string {
	return regexp.QuoteMeta(str)
}

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func ClearBOM(byt []byte) []byte {
	bom := []byte("\ufeff")
	i := bytes.Index([]byte(byt), bom)
	if i == 0 {
		return byt[len(bom):]
	}

	return byt
}

func ReadFileContent(filePath string) (string, error) {
	fileRawContent, err := os.ReadFile(filePath)
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

func GenerateWrapperRegexp(
	wrapperStart string,
	wrapperEnd string,
	wrapperGroup string,
	wrapperLineBreak bool,
) (output string) {
	characters := make([]string, 0, len(wrapperEnd))
	for _, r := range wrapperEnd {
		c := string(r)
		characters = append(characters, c)
	}

	contentPatterns := make([]string, 0)
	if wrapperLineBreak {
		contentPatterns = append(contentPatterns, "\n\\s*")
	}

	for i := 0; i < len(characters); i++ {
		contentPattern := ""
		ci := characters[i]
		for j := 0; j < i; j++ {
			cj := characters[j]
			contentPattern += regexp.QuoteMeta(cj)
		}
		contentPattern += fmt.Sprintf("[^%s]", regexp.QuoteMeta(ci))
		contentPatterns = append(contentPatterns, contentPattern)
	}

	if wrapperLineBreak {
		output = fmt.Sprintf(
			"\\s*%s(\n|\r\n)(?P<%s>(%s)*)(\n|\r\n)\\s*%s",
			regexp.QuoteMeta(wrapperStart),
			wrapperGroup,
			strings.Join(contentPatterns, "|"),
			regexp.QuoteMeta(wrapperEnd),
		)
	} else {
		output = fmt.Sprintf(
			"%s(?P<%s>(%s)*)%s",
			regexp.QuoteMeta(wrapperStart),
			wrapperGroup,
			strings.Join(contentPatterns, "|"),
			regexp.QuoteMeta(wrapperEnd),
		)
	}

	return
}

func IsJsonArray(content string) bool {
	return false
}

func IsJsonObject(content string) bool {
	return false
}
