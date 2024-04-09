package php

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// Serialize turns a Go value into a serialized PHP data string
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
		return serializeString(data.(string)), nil
	case reflect.Struct:
		return serializeStruct(data)
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

func serializeString(value string) string {
	return "s:" + strconv.Itoa(len(value)) + ":\"" + value + "\";"
}

func serializeStruct(data interface{}) (string, error) {
	a := reflect.ValueOf(data)

	var serialized strings.Builder
	serialized.WriteString("O:")
	serialized.WriteString(strconv.Itoa(len(a.Type().Name())))
	serialized.WriteString(":\"")
	serialized.WriteString(a.Type().Name())
	serialized.WriteString("\":")
	serialized.WriteString(strconv.Itoa(structFieldCount(a)))
	serialized.WriteString(":{")

	var field reflect.StructField
	var keyString string
	var valueString string
	var err error
	var ok bool
	for i := range a.NumField() {
		valueString, err = Serialize(a.Field(i).Interface())
		if err != nil {
			return "", err
		}

		field = a.Type().Field(i)
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

	count := 0
	for i := 0; i < data.NumField(); i++ {
		count += structFieldCount(data.Field(i))
	}

	return count
}

// Unserialize populates the destination from a serialized PHP data string
func Unserialize(data string, destination *interface{}) error {
	position := 0
	return unserializeWalk(data, &position, destination)
}

func unserializeWalk(data string, position *int, destination *interface{}) error {
	if data == "N;" {
		return nil
	}

	var err error
	kind := data[*position]

	switch kind {
	case 'a':
		*destination, err = unserializeArray(data, position)
	case 'b':
		*destination, err = unserializeBool(data, position)
	case 'd':
		*destination, err = unserializeFloat(data, position)
	case 'i':
		*destination, err = unserializeInteger(data, position)
	case 's':
		*destination, err = unserializeString(data, position)
	}

	return err
}

func unserializeArray(data string, position *int) (interface{}, error) {
	// Get array size
	*position += 2
	startPosition := *position
	continueUntilCharacter(data, ':', position)
	arrayLength, err := strconv.Atoi(data[startPosition:*position])
	if err != nil {
		return nil, err
	}

	// Skip the ":{"
	*position += 2

	keys := make([]interface{}, arrayLength)
	values := make([]interface{}, arrayLength)

	isAssociative := false
	for i := range arrayLength {
		unserializeWalk(data, position, &keys[i])
		if data[*position] == ';' {
			*position++
		}

		unserializeWalk(data, position, &values[i])
		if data[*position] == ';' {
			*position++
		}

		// Test if array is associative
		if !isAssociative {
			kind := reflect.ValueOf(keys[i]).Kind()
			if kind != reflect.Int || keys[i] != i {
				isAssociative = true
			}
		}
	}

	// Continue until the array is over
	continueUntilCharacter(data, '}', position)
	*position++

	// Consecutive keys detected, create a simple array
	if !isAssociative {
		array := make([]interface{}, 0, arrayLength)
		for i := range arrayLength {
			array = append(array, values[i])
		}

		return array, nil
	}

	array := make(map[interface{}]interface{}, arrayLength)
	for i := range keys {
		array[keys[i]] = values[i]
	}

	return array, nil
}

func unserializeBool(data string, position *int) (bool, error) {
	valueString, err := unserializeValue(data, position)
	if err != nil {
		return false, err
	}

	switch valueString {
	case "0":
		return false, nil
	case "1":
		return true, nil
	default:
		return false, ErrUnexpectedToken{Position: *position - len(valueString)}
	}
}

func unserializeFloat(data string, position *int) (float64, error) {
	valueString, err := unserializeValue(data, position)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(valueString, 64)
}

func unserializeInteger(data string, position *int) (int, error) {
	valueString, err := unserializeValue(data, position)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(valueString)
}

func unserializeString(data string, position *int) (string, error) {
	// Skip string length, maybe later validate this?
	continueUntilCharacter(data, ':', position)
	*position++

	valueString, err := unserializeValue(data, position)
	if err != nil {
		return "", err
	}

	// Validate that we have a nicely quoted string
	if valueString[0] != '"' || valueString[len(valueString)-1] != '"' {
		return valueString, fmt.Errorf("malformed string at position %v", *position-len(valueString))
	}

	return valueString[1 : len(valueString)-1], nil
}

func unserializeValue(data string, position *int) (string, error) {
	continueUntilCharacter(data, ':', position)
	startPosition := *position

	continueUntilCharacter(data, ';', position)

	return data[startPosition+1 : *position], nil
}

func continueUntilCharacter(data string, destination byte, position *int) error {
	for {
		if len(data) == *position {
			return ErrUnexpectedEndOfString
		}

		if data[*position] == destination {
			break
		}

		*position++
	}

	return nil
}
