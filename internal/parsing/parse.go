package parsing

import (
	"log"
	"os"
	"path/filepath"
)

func ParseVariablesFile(path string, jsonPathFilter string, keyColumn string, isMultipleOutput bool, loopInjectionVariable string) (variables []map[string]interface{}, err error) {
	variablesBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	var iVariables interface{}

	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		iVariables, err = ParseJSON(variablesBytes)
		if err != nil {
			log.Fatal(err)
			return
		}
	case ".csv":
		iVariables, err = ParseCSV(variablesBytes, keyColumn)
		if err != nil {
			log.Fatal(err)
			return
		}
	case ".xlsx":
		iVariables, err = ParseXLSX(variablesBytes, keyColumn)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	variables, err = filterAndRootVariables(iVariables, jsonPathFilter, isMultipleOutput, loopInjectionVariable)
	if err != nil {
		return nil, err
	}

	return
}

func ParseVariablesBytes(variablesBytes []byte, jsonPathFilter string, keyColumn string, isMultipleOutput bool, loopInjectionVariable string) (variables []map[string]interface{}, err error) {
	var iVariables interface{}

	iVariables, err = ParseCSV(variablesBytes, keyColumn)
	if err != nil {
		log.Fatal(err)
		return
	}

	// if parameters.MultipleOutput {
	// 	if utils.IsJsonArray(variableContent) {
	// 		variableSet, err = parsing.ParseJSON(variablesContentBytes)
	// 		if err != nil {
	// 			t.SendError(w, 400, err)
	// 			return
	// 		}
	// 	} else {
	// 		variableSet, err = parsing.ParseMultiCSV(variablesContentBytes, parameters.KeyColumn, parameters.WrappingLoopVariable)
	// 		if err != nil {
	// 			t.SendError(w, 400, err)
	// 			return
	// 		}
	// 	}
	// } else {
	// 	variableSet = make([]map[string]interface{}, 1)
	// 	if utils.IsJsonObject(variableContent) {
	// 		variableSet[0], err = parsing.ParseSingleJSON(variablesContentBytes)
	// 		if err != nil {
	// 			t.SendError(w, 400, err)
	// 			return
	// 		}
	// 	} else {
	// 		variableSet[0], err = parsing.ParseSingleCSV(variablesContentBytes, parameters.KeyColumn, parameters.WrappingLoopVariable)
	// 		if err != nil {
	// 			t.SendError(w, 400, err)
	// 			return
	// 		}
	// 	}
	// }

	variables, err = filterAndRootVariables(iVariables, jsonPathFilter, isMultipleOutput, loopInjectionVariable)
	if err != nil {
		return nil, err
	}

	return
}
