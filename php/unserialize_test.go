package php

import (
	"errors"
	"testing"
)

func TestUnserializeError(t *testing.T) {
	err := Unserialize("i:9;", 8)
	if !errors.Is(err, ErrMustBeWriteable) {
		t.Errorf("Expected ErrMustBeWriteable, got %s", err.Error())
	}

	var destInt int
	err = Unserialize("a:1:{i:0;i:1}", &destInt)
	if err == nil {
		t.Error("Expected array error, got none")
	}
}

func TestUnserializeListArray(t *testing.T) {
	var dest []int
	err := Unserialize("a:3:{i:0;i:1;i:1;i:2;i:2;i:3;}", &dest)
	if err != nil {
		t.Error("Unexpected error:", err.Error())
	}
	if len(dest) != 3 || dest[0] != 1 || dest[1] != 2 || dest[2] != 3 {
		t.Errorf("Expected output [1, 2, 3], got %v", dest)
	}

	var aDest [3]int
	err = Unserialize("a:3:{i:0;i:1;i:1;i:2;i:2;i:3;}", &aDest)
	if err != nil {
		t.Error("Unexpected error:", err.Error())
	}
	if aDest[0] != 1 || aDest[1] != 2 || aDest[2] != 3 {
		t.Errorf("Expected output [1, 2, 3], got %v", aDest)
	}

	var destStr []string
	err = Unserialize("a:3:{i:0;s:3:\"foo\";i:1;s:3:\"bar\";i:2;s:3:\"baz\";}", &destStr)
	if err != nil {
		t.Error("Unexpected error:", err.Error())
	}
	if len(destStr) != 3 || destStr[0] != "foo" || destStr[1] != "bar" || destStr[2] != "bar" {
		t.Errorf("Expected output ['foo', 'bar', 'baz'], got %v", destStr)
	}
}

func TestUnserializeObject(t *testing.T) {
	type User struct {
		Name  string
		Admin bool
	}

	var user User
	err := Unserialize("O:4:\"User\":2:{s:4:\"Name\";s:7:\"Jerodev\";s:5:\"Admin\";b:1;}", &user)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}

	if user.Name != "Jerodev" || !user.Admin {
		t.Errorf("Unexpected values %v", user)
	}
}

func TestUnserializeScalar(t *testing.T) {
	var pointer *int
	err := Unserialize("N;", &pointer)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if pointer != nil {
		t.Errorf("expected nil, got %v", pointer)
	}

	var dBool bool
	err = Unserialize("b:1;", &dBool)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if !dBool {
		t.Errorf("expected true, got %v", dBool)
	}

	var dFloat float64
	err = Unserialize("d:8.13;", &dFloat)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dFloat != 8.13 {
		t.Errorf("expected 8.13, got %v", dFloat)
	}

	var dInt int
	err = Unserialize("i:35;", &dInt)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dInt != 35 {
		t.Errorf("expected 35, got %v", dInt)
	}

	err = Unserialize("N;", &dInt)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dInt != 0 {
		t.Errorf("expected 0, got %v", dInt)
	}

	var dString string
	err = Unserialize("s:3:\"foo\";", &dString)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dString != "foo" {
		t.Errorf("expected \"foo\", got %v", dString)
	}
}
