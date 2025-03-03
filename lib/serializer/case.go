package serializer

import (
	"encoding/json"

	"github.com/iancoleman/strcase"
)

func StructToMapUsingSnakeCaseKey(m interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return convertMapKeysToSnake(result), nil
}

// map의 키를 재귀적으로 Snake_Case로 변환한다.
func convertMapKeysToSnake(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		newKey := strcase.ToSnake(k)
		// 만약 값이 map[string]interface{}이면 재귀 호출
		if m, ok := v.(map[string]interface{}); ok {
			v = convertMapKeysToSnake(m)
		} else if arr, ok := v.([]interface{}); ok {
			v = convertSlice(arr)
		}
		out[newKey] = v
	}
	return out
}

// 슬라이스 내부에 map이 있으면 변환
func convertSlice(in []interface{}) []interface{} {
	for i, v := range in {
		if m, ok := v.(map[string]interface{}); ok {
			in[i] = convertMapKeysToSnake(m)
		} else if arr, ok := v.([]interface{}); ok {
			in[i] = convertSlice(arr)
		}
	}
	return in
}
