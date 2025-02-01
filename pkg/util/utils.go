package util

import (
	"fmt"
	"reflect"
)

func GetValueFromPathV2(input *map[string]interface{}, path string) interface{} {
	return (*input)[path]
}

func GetValueFromPath(input map[string]interface{}, path string) interface{} {
	// Simple path resolution for now - could be enhanced to handle nested paths
	return input[path]
}

func CompareValues(a, b interface{}) int {
	// Convert to float64 for numeric comparison
	aFloat, aOk := ToFloat64(a)
	bFloat, bOk := ToFloat64(b)

	if !aOk || !bOk {
		return 0
	}

	if aFloat > bFloat {
		return 1
	} else if aFloat < bFloat {
		return -1
	}
	return 0
}

func ToFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func ContainsValue(container, item interface{}) bool {
	switch c := container.(type) {
	case []interface{}:
		for _, v := range c {
			if reflect.DeepEqual(v, item) {
				return true
			}
		}
	case string:
		itemStr, ok := item.(string)
		if !ok {
			return false
		}
		return Contains(c, itemStr)
	}
	return false
}

func Contains(s string, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-len(substr):] == substr
}

func LogConditionSummary(inputValue interface{}, operationType string, conditionValue interface{}, res bool) {
	fmt.Printf("%v %v %v = %v\n", inputValue, operationType, conditionValue, res)
}
