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
	var rDest reflect.Value
	switch x := dest.(type) {
	case reflect.Value:
		rDest = x
	default:
		rDest = reflect.ValueOf(dest)
		if rDest.Kind() != reflect.Pointer || !rDest.Elem().CanSet() {
			return ErrMustBeWriteable
		}

		rDest = rDest.Elem()
	}

	if input == "N;" {
		rDest.SetZero()
		*position += 3
		return nil
	}

	switch input[*position : *position+2] {
	case "a:":
		return unserializeArray(input, position, dest)
	case "b:":
		if input[*position+2] == '1' {
			rDest.SetBool(true)
		} else {
			rDest.SetBool(false)
		}
		walkUntil(input, position, ';')
	case "d:":
		value, err := strconv.ParseFloat(walkUntil(input, position, ';'), 64)
		if err != nil {
			return err
		}
		rDest.SetFloat(value)
	case "i:":
		value, err := strconv.ParseInt(walkUntil(input, position, ';'), 10, 64)
		if err != nil {
			return err
		}
		rDest.SetInt(value)
	case "O:":
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

func unserializeArray(input string, position *int, dest interface{}) error {
	r := reflect.ValueOf(dest).Elem()
	if r.Kind() == reflect.Map {
		// TODO: implement array => map unserialization

		return errors.New("array to map unserialization not yet implemented")
	}
	if r.Kind() != reflect.Array && r.Kind() != reflect.Slice {
		return errors.New("expected array or slice destination, got " + r.Kind().String())
	}

	return errors.New("array not yet implemented")
}

// unserializeNext finds the next serialized value starting from position and sets the value in destination
func unserializeNext(input string, position *int, dest reflect.Value) error {
	start := *position

	// For simple values, unserialize until the next ;
	if input[*position] != 'O' && input[*position] != 'a' {
		for {
			*position++
			if input[*position] == ';' {
				break
			}
		}

		*position++

		return Unserialize(input[start:*position], dest)
	}

	// TODO: structs and arrays

	return nil
}

func unserializeString(input string, position *int, dest reflect.Value) error {
	*position += 3
	walkUntil(input, position, ':')
	value := walkUntil(input, position, ';')

	if !dest.CanAddr() {
		dest = dest.Elem()
	}

	dest.SetString(value[1 : len(value)-1])
	*position++

	return nil
}

func unserializeStruct(input string, position *int, dest interface{}) error {
	r := reflect.ValueOf(dest).Elem()
	if r.Kind() != reflect.Struct {
		return errors.New("Expected struct destination, got " + r.Kind().String())
	}

	walkUntil(input, position, '{')
	*position++

	var rv reflect.Value
	fieldName := ""
	rFieldName := reflect.ValueOf(&fieldName)
	for {
		if input[*position] == '}' {
			break
		}

		// Find a field name that is part of the struct
		unserializeString(input, position, rFieldName)
		rv = r.FieldByName(fieldName)
		if !rv.IsValid() || !rv.CanSet() {
			return errors.New("Cannot set field " + fieldName)
		}

		// Unserialize the value
		unserializeNext(input, position, rv)
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
