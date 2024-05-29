package php

import (
	"errors"
	"reflect"
	"strconv"
)

var (
	ErrMustBeWriteable = errors.New("destination must be a writeable pointer")
)

func Unserialize(input string, dest interface{}) error {
	position := 0

	return unserializeWalk(input, &position, dest)
}

func unserializeWalk(input string, position *int, dest interface{}) error {
	rDest := reflect.ValueOf(dest)
	if rDest.Kind() != reflect.Pointer || !rDest.Elem().CanSet() {
		return ErrMustBeWriteable
	}

	if input == "N;" {
		rDest.Elem().SetZero()
		*position += 3
		return nil
	}

	switch input[*position : *position+2] {
	case "a:":
		// TODO: unserialize arrays
	case "b:":
		if input[*position+3] == '1' {
			rDest.Elem().SetBool(true)
		} else {
			rDest.Elem().SetBool(false)
		}
		walkUntil(input, position, ';')
	case "d:":
		value, err := strconv.ParseFloat(walkUntil(input, position, ';'), 64)
		if err != nil {
			return err
		}
		rDest.Elem().SetFloat(value)
	case "i:":
		value, err := strconv.Atoi(walkUntil(input, position, ';'))
		if err != nil {
			return err
		}
		rDest.Elem().SetInt(int64(value))
	case "O":
		return unserializeStruct(input, position, dest)
	case "s:":
		return unserializeString(input, position, rDest)
	default:
		return errors.New("unknown pattern " + input)
	}

	// Continue past the last ';'
	*position++

	return nil
}

func unserializeString(input string, position *int, dest reflect.Value) error {
	*position += 3
	walkUntil(input, position, ':')
	value := walkUntil(input, position, ';')
	dest.Elem().SetString(value[1 : len(value)-1])
}

func unserializeStruct(input string, position *int, dest interface{}) error {
	r := reflect.ValueOf(dest)
	if r.Elem().Kind() != reflect.Struct {
		return errors.New("Expected struct destination, got " + r.Elem().Kind().String())
	}

	walkUntil(input, position, '{')
	*position++

	fields := map[string]reflect.Value{}
	for i := range r.Elem().NumField() {
		field := r.Elem().Field(i)
		fields[field.Elem().Type().Name()] = field.Elem()
	}

	fieldName := ""
	rFieldName := reflect.ValueOf(fieldName)
	for {
		if input[*position] == '}' {
			break
		}

		unserializeString(input, position, rFieldName)

		// TODO: map struct field
	}

	*position++

	return nil
}

func walkUntil(input string, position *int, target byte) string {
	startPosition := *position

	for ; *position < len(input); *position++ {
		if input[*position] == target {
			break
		} else if input[*position] == ':' {
			startPosition = *position + 1
		}
	}

	return input[startPosition:*position]
}
