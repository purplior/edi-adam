package dt

import (
	"encoding/json"
	"math"
	"strconv"
)

func Json(in any) (out map[string]interface{}) {
	switch in := in.(type) {
	case nil:
		out = nil
	case int:
		out = nil
	case uint8:
		out = nil
	case uint16:
		out = nil
	case uint32:
		out = nil
	case uint64:
		out = nil
	case uint:
		out = nil
	case int8:
		out = nil
	case int16:
		out = nil
	case int32:
		out = nil
	case int64:
		out = nil
	case float32:
		out = nil
	case float64:
		out = nil
	case bool:
		out = nil
	case string:
		json.Unmarshal([]byte(in), &out)
	case map[string]interface{}:
		out = in
	default:
		b, err := json.Marshal(in)
		if err != nil {
			out = nil
		} else {
			json.Unmarshal(b, &out)
		}
	}

	return out
}

func Int(in any) (out int) {
	out = 0

	switch in := in.(type) {
	case nil:
		out = 0
	case int:
		out = in
	case uint8:
		out = int(in)
	case uint16:
		out = int(in)
	case uint32:
		out = int(in)
	case uint64:
		out = int(in)
	case uint:
		out = int(in)
	case int8:
		out = int(in)
	case int16:
		out = int(in)
	case int32:
		out = int(in)
	case int64:
		out = int(in)
	case float32:
		out = int(math.Floor(float64(in)))
	case float64:
		out = int(math.Floor(in))
	case bool:
		if in {
			out = 1
		}
	case string:
		out, _ = strconv.Atoi(in)
	default:
		out = 1
	}

	return out
}

func Int32(in any) (out int32) {
	out = 0

	switch in := in.(type) {
	case nil:
		out = 0
	case int:
		out = int32(in)
	case uint8:
		out = int32(in)
	case uint16:
		out = int32(in)
	case uint32:
		out = int32(in)
	case uint64:
		out = int32(in)
	case uint:
		out = int32(in)
	case int8:
		out = int32(in)
	case int16:
		out = int32(in)
	case int32:
		out = in
	case int64:
		out = int32(in)
	case float32:
		out = int32(math.Floor(float64(in)))
	case float64:
		out = int32(math.Floor(in))
	case bool:
		if in {
			out = 1
		}
	case string:
		o, _ := strconv.Atoi(in)
		out = int32(o)
	default:
		out = 1
	}

	return out
}

func Int64(in any) (out int64) {
	out = 0

	switch in := in.(type) {
	case nil:
		out = 0
	case int:
		out = int64(in)
	case uint8:
		out = int64(in)
	case uint16:
		out = int64(in)
	case uint32:
		out = int64(in)
	case uint64:
		out = int64(in)
	case uint:
		out = int64(in)
	case int8:
		out = int64(in)
	case int16:
		out = int64(in)
	case int32:
		out = int64(in)
	case int64:
		out = in
	case float32:
		out = int64(math.Floor(float64(in)))
	case float64:
		out = int64(math.Floor(in))
	case bool:
		if in {
			out = 1
		}
	case string:
		o, _ := strconv.Atoi(in)
		out = int64(o)
	default:
		out = 1
	}

	return out
}

func UInt(in any) (out uint) {
	out = 0

	switch in := in.(type) {
	case nil:
		out = 0
	case uint:
		out = in
	case uint8:
		out = uint(in)
	case uint16:
		out = uint(in)
	case uint32:
		out = uint(in)
	case uint64:
		out = uint(in)
	case int:
		out = uint(in)
	case int8:
		out = uint(in)
	case int16:
		out = uint(in)
	case int32:
		out = uint(in)
	case int64:
		out = uint(in)
	case float32:
		out = uint(math.Floor(float64(in)))
	case float64:
		out = uint(math.Floor(in))
	case bool:
		if in {
			out = 1
		}
	case string:
		intOut, _ := strconv.Atoi(in)
		out = UInt(intOut)
	default:
		out = 1
	}

	return out
}

func Float(in any) (out float64) {
	out = 0.0

	switch in := in.(type) {
	case nil:
		out = 0.0
	case int:
		out = float64(in)
	case uint8:
		out = float64(in)
	case uint16:
		out = float64(in)
	case uint32:
		out = float64(in)
	case uint64:
		out = float64(in)
	case uint:
		out = float64(in)
	case int8:
		out = float64(in)
	case int16:
		out = float64(in)
	case int32:
		out = float64(in)
	case int64:
		out = float64(in)
	case float32:
		out = float64(in)
	case float64:
		out = in
	case bool:
		if in {
			out = 1.0
		}
	case string:
		out, _ = strconv.ParseFloat(in, 64)
	default:
		out = 1
	}

	return out
}

func Str(in any) (out string) {
	out = ""

	switch in := in.(type) {
	case nil:
		out = ""
	case int:
		out = strconv.FormatInt(int64(in), 10)
	case uint8:
		out = strconv.FormatInt(int64(in), 10)
	case uint16:
		out = strconv.FormatInt(int64(in), 10)
	case uint32:
		out = strconv.FormatInt(int64(in), 10)
	case uint64:
		out = strconv.FormatInt(int64(in), 10)
	case uint:
		out = strconv.FormatInt(int64(in), 10)
	case int8:
		out = strconv.FormatInt(int64(in), 10)
	case int16:
		out = strconv.FormatInt(int64(in), 10)
	case int32:
		out = strconv.FormatInt(int64(in), 10)
	case int64:
		out = strconv.FormatInt(in, 10)
	case float32:
		out = strconv.FormatFloat(float64(in), 'f', 2, 64)
	case float64:
		out = strconv.FormatFloat(in, 'f', 2, 64)
	case bool:
		out = strconv.FormatBool(in)
	case string:
		out = in
	default:
		out = ""
	}

	return out
}

func Bool(in any) (out bool) {
	out = false

	switch in := in.(type) {
	case nil:
		out = false
	case int:
		out = in > 0
	case uint8:
		out = in > 0
	case uint16:
		out = in > 0
	case uint32:
		out = in > 0
	case uint64:
		out = in > 0
	case uint:
		out = in > 0
	case int8:
		out = in > 0
	case int16:
		out = in > 0
	case int32:
		out = in > 0
	case int64:
		out = in > 0
	case float32:
		out = in > 0
	case float64:
		out = in > 0
	case bool:
		out = in
	case string:
		out = len(in) > 0
	default:
		out = true
	}

	return out
}
