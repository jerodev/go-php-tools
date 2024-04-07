package php

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func Serialize(data interface{}) (string, error) {
	if data == nil {
		return "N;", nil
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return serializeArray(data)
	case reflect.Bool:
		if data.(bool) {
			return "b:1;", nil
		}
		return "b:0;", nil
	case reflect.Float32:
	case reflect.Float64:
		return "d:" + fmt.Sprintf("%v", data) + ";", nil
	case reflect.Int:
		return "i:" + strconv.Itoa(data.(int)) + ";", nil
	case reflect.Int8:
		return "i:" + strconv.Itoa(int(data.(int8))) + ";", nil
	case reflect.Int16:
		return "i:" + strconv.Itoa(int(data.(int16))) + ";", nil
	case reflect.Int32:
		return "i:" + strconv.Itoa(int(data.(int32))) + ";", nil
	case reflect.Int64:
		return "i:" + strconv.Itoa(int(data.(int64))) + ";", nil
	case reflect.Map:
		return serializeMap(data)
	case reflect.String:
		return "s:" + strconv.Itoa(len(data.(string))) + ":\"" + data.(string) + "\";", nil
	case reflect.Uint:
		return "i:" + strconv.Itoa(int(data.(uint))) + ";", nil
	case reflect.Uint8:
		return "i:" + strconv.Itoa(int(data.(uint8))) + ";", nil
	case reflect.Uint16:
		return "i:" + strconv.Itoa(int(data.(uint16))) + ";", nil
	case reflect.Uint32:
		return "i:" + strconv.Itoa(int(data.(uint32))) + ";", nil
	case reflect.Uint64:
		return "i:" + strconv.Itoa(int(data.(uint64))) + ";", nil
	default:
		return "", ErrUnsupportedDataType{DataType: v.Kind().String()}
	}

	return "", ErrUnsupportedDataType{DataType: v.Kind().String()}
}

func serializeArray(data interface{}) (string, error) {
	a := reflect.ValueOf(data)

	var serialized strings.Builder
	serialized.WriteString("a:" + strconv.Itoa(a.Len()) + ":{")

	var valueString string
	var err error
	for i := 0; i < a.Len(); i++ {
		valueString, err = Serialize(reflect.Indirect(a.Index(i)).Interface())
		if err != nil {
			return "", err
		}

		serialized.WriteString("i:")
		serialized.WriteString(strconv.Itoa(i))
		serialized.WriteString(";")
		serialized.WriteString(valueString)
	}

	return serialized.String() + "}", nil
}

func serializeMap(data interface{}) (string, error) {
	a := reflect.ValueOf(data)

	var serialized strings.Builder
	serialized.WriteString("a:" + strconv.Itoa(a.Len()) + ":{")

	// The order of keys in a map is unpredictable.
	// We sort the keys alphabetically to make testing easier
	keys := a.MapKeys()
	slices.SortFunc(keys, func(a, b reflect.Value) int {
		if a.Kind() == reflect.String {
			return strings.Compare(a.String(), b.String())
		} else {
			return int(a.Int() - b.Int())
		}
	})

	var keyString string
	var valueString string
	var err error
	for _, k := range keys {
		keyString, err = Serialize(k.Interface())
		if err != nil {
			return "", err
		}

		valueString, err = Serialize(a.MapIndex(k).Interface())
		if err != nil {
			return "", err
		}

		serialized.WriteString(strings.TrimRight(keyString, ";"))
		serialized.WriteString(";")
		serialized.WriteString(valueString)
	}

	return serialized.String() + "}", nil
}
