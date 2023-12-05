package ureflect

import "reflect"

func ConvertInterfaceToSlice(input interface{}) ([]interface{}, bool) {
	// 使用反射获取input的值
	value := reflect.ValueOf(input)

	// 检查input是否是切片类型
	if value.Kind() != reflect.Slice {
		return nil, false
	}

	length := value.Len()
	result := make([]interface{}, length)

	// 将切片中的元素转换为interface{}类型并存储到result中
	for i := 0; i < length; i++ {
		result[i] = value.Index(i).Interface()
	}

	return result, true
}

func IsSlice(input interface{}) bool {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Slice {
		return true
	}

	return false
}
