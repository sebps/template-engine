package parsing

import (
	"errors"
	"log"

	"github.com/sebps/template-engine/internal/filtering"
	"github.com/sebps/template-engine/internal/utils"
)

func filterVariables(input interface{}, jsonPathFilter string) (output interface{}, err error) {
	filtered := input
	if len(jsonPathFilter) > 0 {
		filtered = filtering.Filter(input, jsonPathFilter)
	}

	if utils.IsArray(filtered) {
		var arrayOutput []map[string]interface{}
		utils.MarshalUnmarshal(filtered, &arrayOutput)
		output = arrayOutput
	} else {
		var mapOutput map[string]interface{}
		utils.MarshalUnmarshal(filtered, &mapOutput)
		output = mapOutput
	}

	return
}

func filterAndRootVariables(iVariables interface{}, jsonPathFilter string, isMultipleOutput bool, loopInjectionVariable string) (variables []map[string]interface{}, err error) {
	fVariables, err := filterVariables(iVariables, jsonPathFilter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if utils.IsArray(fVariables) && isMultipleOutput {
		variables = fVariables.([]map[string]interface{})
	} else if utils.IsArray(fVariables) && !isMultipleOutput {
		// root flat variable to prepare for template injection
		rootVariables, err := RootVariables(fVariables, loopInjectionVariable)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		variables = make([]map[string]interface{}, 1)
		variables[0] = rootVariables
	} else if !utils.IsArray(fVariables) && isMultipleOutput {
		err = errors.New("multiple output requires a flat data input")
	} else if !utils.IsArray(fVariables) && !isMultipleOutput {
		variables = make([]map[string]interface{}, 1)
		variables[0] = fVariables.(map[string]interface{})
	}

	return
}
