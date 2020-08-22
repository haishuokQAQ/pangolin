package utils

import (
	"errors"
	"reflect"
	"strconv"
)

func StructToQueryMap(value interface{}) (map[string]string, error) {
	result := map[string]string{}
	obj := reflect.ValueOf(value)
	elem := obj.Elem()
	if elem.Kind() == reflect.Struct {
		elemType := elem.Type()
		for i := 0; i < elemType.NumField(); i++ {
			typeStruct := elemType.Field(i)
			fieldKind := elem.Field(i).Kind()
			fieldValue := elem.Field(i).Interface()
			var valueString string
			var valueInt int64
			isInt := true
			switch fieldKind {
			case reflect.Int:
				actualValue := fieldValue.(int)
				valueInt = int64(actualValue)
				break
			case reflect.Int8:
				actualValue := fieldValue.(int8)
				valueInt = int64(actualValue)
				break
			case reflect.Int16:
				actualValue := fieldValue.(int16)
				valueInt = int64(actualValue)
				break
			case reflect.Int32:
				actualValue := fieldValue.(int32)
				valueInt = int64(actualValue)
				break
			case reflect.Int64:
				valueInt = fieldValue.(int64)
				break
			default:
				isInt = false
			}
			var valueUint uint64
			isUint := true
			switch fieldKind {
			case reflect.Int:
				actualValue := fieldValue.(int)
				valueInt = int64(actualValue)
				break
			case reflect.Int8:
				actualValue := fieldValue.(int8)
				valueInt = int64(actualValue)
				break
			case reflect.Int16:
				actualValue := fieldValue.(int16)
				valueInt = int64(actualValue)
				break
			case reflect.Int32:
				actualValue := fieldValue.(int32)
				valueInt = int64(actualValue)
				break
			case reflect.Int64:
				valueInt = fieldValue.(int64)
				break
			default:
				isInt = false
			}
			if isUint {
				valueString = strconv.FormatUint(valueUint, 10)
			} else if isInt {
				valueString = strconv.FormatInt(valueInt, 10)
			}

			if fieldKind == reflect.Int ||
				fieldKind == reflect.Int8 ||
				fieldKind == reflect.Int16 ||
				fieldKind == reflect.Int32 ||
				fieldKind == reflect.Int64 {

			}
			result[typeStruct.Name] = valueString
		}
	} else {
		return nil, errors.New("Need struct converting to query map!")
	}
	return result, nil
}
