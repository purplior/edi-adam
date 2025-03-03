package dt_test

import (
	"testing"

	"github.com/purplior/edi-adam/lib/dt"
)

func Test_Json_with_string_argument(t *testing.T) {
	given := "{\"a\": 1}"

	data := dt.Json(given)
	value, has := data["a"]
	success := has && value == 1.0

	if !success {
		t.Logf("value: %.0f, has: %t", value, has)
		t.Fail()
	}
}

func Test_Json_with_struct_argument(t *testing.T) {
	given := struct {
		A string  `json:"a"`
		B float64 `json:"b"`
	}{
		A: "hello",
		B: 3.0,
	}

	data := dt.Json(given)
	valueA, hasA := data["a"]
	valueB, hasB := data["b"]
	success := hasA && hasB && valueA == "hello" && valueB == 3.0

	if !success {
		t.Logf("valueA: %s, hasA: %t, valueB: %.0f, hasB: %t", valueA, hasA, valueB, hasB)
		t.Fail()
	}
}
