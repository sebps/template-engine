package filtering

import (
	"github.com/sebps/jsonpath"
)

func Filter(input interface{}, jsonPathFilter string) (output interface{}) {
	filteredVariables, err := jsonpath.JsonPathLookup(input, jsonPathFilter)
	if err != nil {
		panic("could not filter data based on jsonpath query")
	}

	return filteredVariables
}
