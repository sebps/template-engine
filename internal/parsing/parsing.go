package parsing

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebps/template-engine/internal/utils"
)

func parseJSON(path string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	variablesBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func parseCSV(path string, keyCol string, loopVariable string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fileReader := csv.NewReader(fd)
	fileReader.Comma = ';'
	fileReader.LazyQuotes = true

	records := make(map[string][]string)
	cols := make(map[int]string)
	rows := make(map[int]string)
	orderedCols := make([]string, 0)

	rowNum := 0
	for {
		record, err := fileReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if rowNum == 0 {
			for colNum, col := range record {
				colName := utils.ClearString(col)
				records[colName] = make([]string, 0)
				cols[colNum] = colName
				orderedCols = append(orderedCols, colName)
			}
		} else {
			for colNum, value := range record {
				col := cols[colNum]
				if strings.Compare(utils.ClearString(col), keyCol) == 0 {
					rows[rowNum-1] = strings.TrimSpace(value)
				} else {
					records[col] = append(records[col], value)
				}
			}
		}

		rowNum++
	}

	formattedVariables := make(map[string][]map[string]string)
	formattedVariables[loopVariable] = make([]map[string]string, len(records)-1)
	currentIndex := 0

	for _, colName := range orderedCols {
		if strings.Compare(colName, keyCol) == 0 {
			continue
		}

		recordValues := records[colName]
		formattedVariables[loopVariable][currentIndex] = make(map[string]string)
		for sliceIndex, sliceValue := range recordValues {
			variable := rows[sliceIndex]
			if variable != "" {
				formattedVariables[loopVariable][currentIndex][variable] = sliceValue
			}
		}
		currentIndex++
	}

	variablesBytes, err := json.Marshal(formattedVariables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func ParseVariables(path string, keyColumn string, loopVariable string) map[string]interface{} {
	var variables map[string]interface{}
	var err error

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		variables, err = parseJSON(path)
		if err != nil {
			log.Fatal(err)
		}
	case ".csv":
		variables, err = parseCSV(path, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	}

	return variables
}
