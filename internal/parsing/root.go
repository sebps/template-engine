package parsing

import "encoding/json"

func RootVariables(input interface{}, root string) (output map[string]interface{}, err error) {
	output = make(map[string]interface{})
	bInput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var iInput interface{}
	err = json.Unmarshal(bInput, &iInput)
	if err != nil {
		return nil, err
	}
	output[root] = iInput

	return
}
