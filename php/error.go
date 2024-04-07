package php

import "fmt"

type ErrUnsupportedDataType struct {
	DataType string
}

func (e ErrUnsupportedDataType) Error() string {
	return fmt.Sprintf("Unsupported data type %s", e.DataType)
}
