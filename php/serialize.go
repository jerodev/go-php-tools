package php

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var structNames map[string]string

// Serialize turns a Go value into a serialized PHP data string
func Serialize(data interface{}) (string, error) {
	if data == nil {
		return "N;", nil
	}

	if structNames == nil {
		structNames = map[string]string{}
	}

	var v reflect.Value
	var ok bool
	if v, ok = data.(reflect.Value); !ok {
		v = reflect.ValueOf(data)
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return serializeArray(v)
	case reflect.Bool:
		if data.(bool) {
			return "b:1;", nil
		}
		return "b:0;", nil
	case reflect.Float32:
	case reflect.Float64:
		return "d:" + fmt.Sprintf("%v", v.Float()) + ";", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "i:" + strconv.Itoa(int(v.Int())) + ";", nil
	case reflect.Map:
		return serializeMap(v)
	case reflect.String:
		return serializeString(v.String()), nil
	case reflect.Struct:
		return serializeStruct(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "i:" + strconv.Itoa(int(v.Uint())) + ";", nil
	case reflect.Pointer:
		return Serialize(v.Elem())
	default:
		return "", ErrUnsupportedDataType{DataType: v.Kind().String()}
	}

	return "", ErrUnsupportedDataType{DataType: v.Kind().String()}
}

// WithStructNames allows setting the desired object name for a Go struct
// The keys in the map represent the Go struct names and the values are the PHP class names
func WithStructNames(names map[string]string) {
	if structNames == nil {
		structNames = names
		return
	}

	for k, v := range names {
		if v == "" {
			delete(structNames, k)
		} else {
			structNames[k] = v
		}
	}
}

func serializeArray(data reflect.Value) (string, error) {
	var serialized strings.Builder
	serialized.WriteString("a:" + strconv.Itoa(data.Len()) + ":{")

	var valueString string
	var err error
	for i := range data.Len() {
		valueString, err = Serialize(reflect.Indirect(data.Index(i)))
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

func serializeMap(data reflect.Value) (string, error) {
	var serialized strings.Builder
	serialized.WriteString("a:" + strconv.Itoa(data.Len()) + ":{")

	// The order of keys in a map is unpredictable.
	// We sort the keys alphabetically to make the output consistent
	// Take into account key casting: https://www.php.net/manual/en/language.types.array.php#language.types.array.key-casts
	keys := data.MapKeys()
	slices.SortFunc(keys, func(a, b reflect.Value) int {
		if a.CanInt() {
			return int(a.Int() - b.Int())
		}
		if a.Kind() == reflect.String {
			return strings.Compare(a.String(), b.String())
		}
		if a.CanFloat() {
			return int(a.Float()) - int(b.Float())
		}
		if a.Kind() == reflect.Bool {
			if a.Bool() && !b.Bool() {
				return 1
			}
			if !a.Bool() && b.Bool() {
				return -1
			}
			return 0
		}

		return -1
	})

	var keyString string
	var valueString string
	var keyValue interface{}
	var err error
	for _, k := range keys {
		if k.Kind() == reflect.Struct {
			return "", errors.New("Arrays and objects can not be used as keys. Received " + k.Kind().String())
		}

		// Cast array keys according to PHP specs
		// https://www.php.net/manual/en/language.types.array.php#language.types.array.key-casts
		keyValue = k.Interface()
		if k.CanFloat() {
			keyValue = int(k.Float())
		} else if keyValue == nil || (k.Kind() == reflect.Pointer && k.IsNil()) {
			keyValue = ""
		} else if k.Kind() == reflect.Bool {
			if k.Bool() {
				keyValue = 1
			} else {
				keyValue = 0
			}
		}

		keyString, err = Serialize(keyValue)
		if err != nil {
			return "", err
		}

		valueString, err = Serialize(data.MapIndex(k))
		if err != nil {
			return "", err
		}

		serialized.WriteString(strings.TrimSuffix(keyString, ";"))
		serialized.WriteString(";")
		serialized.WriteString(valueString)
	}

	return serialized.String() + "}", nil
}

func serializeString(value string) string {
	return "s:" + strconv.Itoa(len(value)) + ":\"" + value + "\";"
}

func serializeStruct(data reflect.Value) (string, error) {
	objectName := data.Type().Name()
	if name, ok := structNames[objectName]; ok {
		objectName = name
	}

	var serialized strings.Builder
	serialized.WriteString("O:")
	serialized.WriteString(strconv.Itoa(len(objectName)))
	serialized.WriteString(":\"")
	serialized.WriteString(objectName)
	serialized.WriteString("\":")
	serialized.WriteString(strconv.Itoa(structFieldCount(data)))
	serialized.WriteString(":{")

	var field reflect.StructField
	var keyString, valueString string
	var err error
	var ok bool
	for i := range data.NumField() {
		valueString, err = Serialize(data.Field(i).Interface())
		if err != nil {
			return "", err
		}

		field = data.Type().Field(i)
		if keyString, ok = field.Tag.Lookup("php"); !ok {
			keyString = field.Name
		}
		serialized.WriteString(serializeString(keyString))

		serialized.WriteString(valueString)
	}

	serialized.WriteByte('}')

	return serialized.String(), nil
}

func structFieldCount(data reflect.Value) int {
	if data.Kind() != reflect.Struct || data.IsZero() {
		return 1
	}

	var count int
	for i := range data.NumField() {
		count += structFieldCount(data.Field(i))
	}

	return count
}
