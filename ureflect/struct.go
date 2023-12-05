package ureflect

import (
	"fmt"
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}

	return data
}

func GetMethodNames(obj interface{}) []string {
	objType := reflect.TypeOf(obj)
	numMethods := objType.NumMethod()
	methodNames := make([]string, 0, numMethods)
	for i := 0; i < numMethods; i++ {
		methodName := objType.Method(i).Name
		methodNames = append(methodNames, methodName)
	}

	return methodNames
}

func GetMethodMap(obj interface{}) ([]map[string]interface{}, error) {
	objType := reflect.TypeOf(obj)

	numMethods := objType.NumMethod()

	methodMaps := make([]map[string]interface{}, 0)

	for i := 0; i < numMethods; i++ {
		method := objType.Method(i)
		methodName := method.Name
		fmt.Println(method.PkgPath)
		methodValue := method.Func.Interface()

		methodMap := map[string]interface{}{
			"name":  methodName,
			"value": methodValue,
		}
		methodMaps = append(methodMaps, methodMap)
	}

	return methodMaps, nil
}

func GetPackageName(obj interface{}) string {
	objType := reflect.TypeOf(obj)
	return objType.PkgPath()
}

func GetStructName(obj interface{}) string {
	objType := reflect.TypeOf(obj)

	if objType.Kind() == reflect.Struct {
		return objType.Name()
	} else {
		return "Not a struct"
	}
}
