package filtering

import (
	"encoding/json"

	"github.com/sebps/jsonpath"
)

func Filter(input interface{}, jsonPathFilter string) (output interface{}) {
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
