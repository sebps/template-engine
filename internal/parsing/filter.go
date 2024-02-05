package parsing

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/sebps/template-engine/internal/filtering"
	"github.com/sebps/template-engine/internal/utils"
)

func filterVariables(input interface{}, jsonPathFilter string) (output []map[string]interface{}, err error) {
	var tmp interface{} = input

	if len(jsonPathFilter) > 0 {
		tmp = filtering.Filter(tmp, jsonPathFilter)
	}

	bTmp, err := json.Marshal(tmp)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	json.Unmarshal(bTmp, &output)

	return
}

func filterAndRootVariables(iVariables interface{}, jsonPathFilter string, isMultipleOutput bool, loopInjectionVariable string) (variables []map[string]interface{}, err error) {
	fVariables, err := filterVariables(iVariables, jsonPathFilter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	variables = fVariables
	if isMultipleOutput && !utils.IsArray(fVariables) {
		err = errors.New("multiple output requires list data type")
	} else if !isMultipleOutput && utils.IsArray(fVariables) {
		rootVariables, err := RootVariables(fVariables, loopInjectionVariable)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		variables = make([]map[string]interface{}, 1)
		variables[0] = rootVariables
	}

	return
}
