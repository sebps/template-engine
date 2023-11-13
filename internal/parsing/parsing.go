package parsing

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebps/template-engine/internal/utils"
)

func parseMultiJSON(path string) ([]map[string]interface{}, error) {
	var variables []map[string]interface{}

	variablesBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func parseSingleJSON(path string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	variablesBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func doParseCSV(path string, keyCol string, loopVariable string) ([]map[string]string, error) {
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
				if colNum == len(record)-1 && col == "" {
					continue
				}
				colName := utils.ClearString(col)
				records[colName] = make([]string, 0)
				cols[colNum] = colName
				orderedCols = append(orderedCols, colName)
			}
		} else {
			for colNum, value := range record {
				col := cols[colNum]
				colName := utils.ClearString(col)
				if strings.Compare(colName, keyCol) == 0 {
					rows[rowNum-1] = strings.TrimSpace(value)
				} else {
					records[col] = append(records[col], value)
				}
			}
		}

		rowNum++
	}

	rootLoop := make([]map[string]string, len(orderedCols)-1)
	for currentIndex, colName := range orderedCols {
		if strings.Compare(colName, keyCol) == 0 {
			continue
		}

		recordValues := records[colName]
		rootLoop[currentIndex] = make(map[string]string)
		for sliceIndex, sliceValue := range recordValues {
			variable := rows[sliceIndex]
			if variable != "" {
				rootLoop[currentIndex][variable] = sliceValue
			}
		}
	}

	return rootLoop, nil
}

func parseSingleCSV(path string, keyCol string, loopVariable string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	rootLoop, err := doParseCSV(path, keyCol, loopVariable)
	if err != nil {
		return nil, err
	}

	formattedVariables := make(map[string][]map[string]string)
	formattedVariables[loopVariable] = rootLoop

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

func parseMultiCSV(path string, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
	var variables []map[string]interface{}

	formattedVariables, err := doParseCSV(path, keyCol, loopVariable)
	if err != nil {
		return nil, err
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

func ParseSingleVariables(path string, keyColumn string, loopVariable string) map[string]interface{} {
	var variables map[string]interface{}
	var err error

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		variables, err = parseSingleJSON(path)
		if err != nil {
			log.Fatal(err)
		}
	case ".csv":
		variables, err = parseSingleCSV(path, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	}

	return variables
}

func ParseMultiVariables(path string, keyColumn string, loopVariable string) []map[string]interface{} {
	var variables []map[string]interface{}
	var err error

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		variables, err = parseMultiJSON(path)
		if err != nil {
			log.Fatal(err)
		}
	case ".csv":
		variables, err = parseMultiCSV(path, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	}

	return variables
}
