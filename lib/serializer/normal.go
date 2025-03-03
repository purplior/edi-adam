package serializer

import "encoding/json"

func MapToStruct(m map[string]interface{}, dest interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
