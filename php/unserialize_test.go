package php

import "testing"

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
	err = Unserialize("b:8.13;", &dBool)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dBool {
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

	var dString string
	err = Unserialize("s:3:\"foo\";", &dString)
	if err != nil {
		t.Errorf("ERR %s", err.Error())
	}
	if dString != "foo" {
		t.Errorf("expected \"foo\", got %v", dString)
	}
}
