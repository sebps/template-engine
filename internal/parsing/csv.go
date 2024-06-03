package parsing

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"

	"github.com/sebps/template-engine/internal/utils"
)

func doParseCSV(data []byte, keyCol string) ([]map[string]interface{}, error) {
	r := bytes.NewReader(data)

	fileReader := csv.NewReader(r)
	fileReader.Comma = ';'
	fileReader.LazyQuotes = true

	rootLoop := make([]map[string]interface{}, 0)

	keyColNum := -1
	rowNum := 0
	for {
		row, err := fileReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if rowNum == 0 {
			for colNum, colName := range row {
				if colName == keyCol {
					keyColNum = colNum
				} else {
					rootLoop = append(rootLoop, make(map[string]interface{}))
				}
			}
			if keyColNum == -1 {
				return nil, errors.New("key column not found")
			}
		} else {
			if len(row)-1 < keyColNum {
				// if row length too small to reach key col index no further processing of the current row
				continue
			}
			currentVariable := row[keyColNum]
			if currentVariable == "" {
				// if no variable no further processing of the current row
				continue
			}
			for colNum, colValue := range row {
				if colNum < keyColNum {
					rootLoop[colNum][currentVariable] = colValue
				} else if colNum > keyColNum {
					rootLoop[colNum-1][currentVariable] = colValue
				}
			}
		}

		rowNum++
	}

	return rootLoop, nil
}

func ParseCSV(data []byte, keyCol string) (variables interface{}, err error) {
	data = utils.ClearBOM(data)

	rootLoop, err := doParseCSV(data, keyCol)
	if err != nil {
		return nil, err
	}

	variablesBytes, err := json.Marshal(rootLoop)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return
}
