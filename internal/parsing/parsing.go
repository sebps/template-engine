package parsing

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebps/jsonpath"
	"github.com/sebps/template-engine/internal/utils"
	"github.com/xuri/excelize/v2"
)

func ParseMultiJSON(data []byte) ([]map[string]interface{}, error) {
	var variables []map[string]interface{}
	err := json.Unmarshal(data, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func ParseSingleJSON(variablesBytes []byte) (map[string]interface{}, error) {
	var variables map[string]interface{}
	err := json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func doParseCSV(data []byte, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
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

func doParseCSV_v0(data []byte, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
	r := bytes.NewReader(data)

	fileReader := csv.NewReader(r)
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

	rootLoop := make([]map[string]interface{}, len(orderedCols)-1)
	for currentIndex, colName := range orderedCols {
		if strings.Compare(colName, keyCol) == 0 {
			continue
		}

		recordValues := records[colName]
		rootLoop[currentIndex] = make(map[string]interface{})
		for sliceIndex, sliceValue := range recordValues {
			variable := rows[sliceIndex]
			if variable != "" {
				rootLoop[currentIndex][variable] = sliceValue
			}
		}
	}

	return rootLoop, nil
}

func doParseXLSX(data []byte, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
	r := bytes.NewReader(data)

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// process first sheet only
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	rootLoop := make([]map[string]interface{}, 0)

	keyColNum := -1
	rowNum := 0
	for _, row := range rows {
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

func ParseSingleCSV(data []byte, keyCol string, loopVariable string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	rootLoop, err := doParseCSV(data, keyCol, loopVariable)
	if err != nil {
		return nil, err
	}

	formattedVariables := make(map[string][]map[string]interface{})
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

func ParseMultiCSV(data []byte, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
	var variables []map[string]interface{}

	formattedVariables, err := doParseCSV(data, keyCol, loopVariable)
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

func ParseSingleXLSX(data []byte, keyCol string, loopVariable string) (map[string]interface{}, error) {
	var variables map[string]interface{}

	rootLoop, err := doParseXLSX(data, keyCol, loopVariable)
	if err != nil {
		return nil, err
	}

	formattedVariables := make(map[string][]map[string]interface{})
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

func ParseMultiXLSX(data []byte, keyCol string, loopVariable string) ([]map[string]interface{}, error) {
	var variables []map[string]interface{}

	formattedVariables, err := doParseXLSX(data, keyCol, loopVariable)
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

func ParseSingleVariablesFile(path string, keyColumn string, loopVariable string, jsonPathFilter string) map[string]interface{} {
	var variables map[string]interface{}
	var err error

	variablesBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		variables, err = ParseSingleJSON(variablesBytes)
		if err != nil {
			log.Fatal(err)
		}
	case ".csv":
		variables, err = ParseSingleCSV(variablesBytes, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	case ".xlsx":
		variables, err = ParseSingleXLSX(variablesBytes, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	}

	variables = filterSingleVariables(variables, jsonPathFilter)

	return variables
}

func ParseMultiVariablesFile(path string, keyColumn string, loopVariable string, jsonPathFilter string) []map[string]interface{} {
	var variables []map[string]interface{}
	var err error

	variablesBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		variables, err = ParseMultiJSON(variablesBytes)
		if err != nil {
			log.Fatal(err)
		}
	case ".csv":
		variables, err = ParseMultiCSV(variablesBytes, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	case ".xlsx":
		variables, err = ParseMultiXLSX(variablesBytes, keyColumn, loopVariable)
		if err != nil {
			log.Fatal(err)
		}
	}

	variables = filterMultiVariables(variables, jsonPathFilter)

	return variables
}

func filterSingleVariables(input map[string]interface{}, jsonPathFilter string) (output map[string]interface{}) {
	output = make(map[string]interface{})

	filteredVariables, err := jsonpath.JsonPathLookup(input, jsonPathFilter)
	if err != nil {
		panic("could not filter data based on jsonpath query")
	}

	if len(filteredVariables.([]interface{})) > 0 {
		marshalled, err := json.Marshal(filteredVariables.([]interface{})[0])
		if err != nil {
			panic("could not filter data based on jsonpath query")
		}

		json.Unmarshal(marshalled, &output)
	} else {
		panic("could not filter data based on jsonpath query")
	}

	return
}

func filterMultiVariables(input []map[string]interface{}, jsonPathFilter string) (output []map[string]interface{}) {
	output = make([]map[string]interface{}, 0)

	filteredVariables, err := jsonpath.JsonPathLookup(input, jsonPathFilter)
	if err != nil {
		panic("could not filter data based on jsonpath query")
	}

	if len(filteredVariables.([]interface{})) > 0 {
		marshalled, err := json.Marshal(filteredVariables)
		if err != nil {
			panic("could not filter data based on jsonpath query")
		}

		json.Unmarshal(marshalled, &output)
	}

	return
}
