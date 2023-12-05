package uslice

import (
	"git.umu.work/AI/uglib/ureflect"
)

type Item struct {
	Value interface{}
}

func IsSlice(slice interface{}) bool {
	switch slice.(type) {
	case []interface{}, []*Item, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64,
		[]float32, []float64, []string, []bool:
		return true
	default:
		ok := ureflect.IsSlice(slice)
		return ok
	}
}

func Convert2Slice(slice interface{}) ([]interface{}, bool) {
	switch s := slice.(type) {
	case []interface{}:
		return s, true
	case []*Item:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v.Value
		}
		return result, true
	case []int:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []int8:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []int16:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []int32:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []int64:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []uint:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []uint8:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []uint16:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []uint32:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []uint64:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []float32:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []float64:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []string:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	case []bool:
		result := make([]interface{}, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result, true
	default:
		result, ok := ureflect.ConvertInterfaceToSlice(slice)
		return result, ok
	}
}

func FlattenNestedSlice(slice interface{}) ([]interface{}, bool) {
	result, ok := Convert2Slice(slice)
	if !ok {
		return nil, false
	}
	flattened := make([]interface{}, 0, len(result))
	for _, item := range result {
		ok := IsSlice(item)
		if ok {
			// 递归调用扁平化函数
			nestedFlattened, _ := FlattenNestedSlice(item)
			flattened = append(flattened, nestedFlattened...)
		} else {
			flattened = append(flattened, item)
		}
	}

	return flattened, true
}
