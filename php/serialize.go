package php

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Serialize(data interface{}) string {
	if data == nil {
		return "N;"
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return serializeArray(data)
	case reflect.Bool:
		if data.(bool) {
			return "b:1;"
		}
		return "b:0;"
	case reflect.Float32:
	case reflect.Float64:
		return "d:" + fmt.Sprintf("%v", data) + ";"
	case reflect.Int:
		return "i:" + strconv.Itoa(data.(int)) + ";"
	case reflect.Int8:
		return "i:" + strconv.Itoa(int(data.(int8))) + ";"
	case reflect.Int16:
		return "i:" + strconv.Itoa(int(data.(int16))) + ";"
	case reflect.Int32:
		return "i:" + strconv.Itoa(int(data.(int32))) + ";"
	case reflect.Int64:
		return "i:" + strconv.Itoa(int(data.(int64))) + ";"
	case reflect.String:
		return "s:" + strconv.Itoa(len(data.(string))) + ":\"" + data.(string) + "\";"
	case reflect.Uint:
		return "i:" + strconv.Itoa(int(data.(uint))) + ";"
	case reflect.Uint8:
		return "i:" + strconv.Itoa(int(data.(uint8))) + ";"
	case reflect.Uint16:
		return "i:" + strconv.Itoa(int(data.(uint16))) + ";"
	case reflect.Uint32:
		return "i:" + strconv.Itoa(int(data.(uint32))) + ";"
	case reflect.Uint64:
		return "i:" + strconv.Itoa(int(data.(uint64))) + ";"
	default:
		return ""
	}

	return ""
}

func serializeArray(data interface{}) string {
	a := reflect.ValueOf(data)

	var serialized strings.Builder
	serialized.WriteString("a:" + strconv.Itoa(a.Len()) + ":{")

	for i := 0; i < a.Len(); i++ {
		v := reflect.Indirect(a.Index(i))

		serialized.WriteString("i:")
		serialized.WriteString(strconv.Itoa(i))
		serialized.WriteString(";")
		serialized.WriteString(Serialize(v.Interface()))
	}

	return serialized.String() + "}"
}
