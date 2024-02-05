package parsing

import (
	"encoding/json"
)

func ParseJSON(variablesBytes []byte) (variables interface{}, err error) {
	err = json.Unmarshal(variablesBytes, &variables)
	if err != nil {
		return nil, err
	}

	return
}
