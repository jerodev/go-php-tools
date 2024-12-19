package php

import (
	"fmt"
)

type ErrUnexpectedToken struct {
	Position int
}

func (e ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token at position %v", e.Position)
}

type ErrUnsupportedDataType struct {
	DataType string
}

func (e ErrUnsupportedDataType) Error() string {
	return fmt.Sprintf("unsupported data type %s", e.DataType)
}
