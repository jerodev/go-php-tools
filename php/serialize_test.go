package php

import "testing"

func TestSerializeArray(t *testing.T) {
	expected := "a:3:{i:0;i:4;i:1;i:5;i:2;i:6;}"
	if expected != Serialize([3]int{4, 5, 6}) {
		t.Errorf("output did not match expecation, got %s", Serialize([3]int{4, 5, 6}))
	}

	expected = "a:2:{i:0;i:1;i:1;i:2;}"
	if expected != Serialize([]int{1, 2}) {
		t.Errorf("output did not match expecation, got %s", Serialize([]int{1, 2}))
	}

	expected = "a:3:{i:0;a:2:{i:0;i:1;i:1;i:2;}i:1;a:2:{i:0;i:3;i:1;i:4;}i:2;a:1:{i:0;i:5;}}"
	if expected != Serialize([][]int{{1, 2}, {3, 4}, {5}}) {
		t.Errorf("output did not match expecation, got %s", Serialize([][]int{{1, 2}, {3, 4}, {5}}))
	}
}

func TestSerializeScalar(t *testing.T) {
	trials := map[any]string{
		nil:       "N;",
		true:      "b:1;",
		false:     "b:0;",
		1.23:      "d:1.23;",
		123:       "i:123;",
		"foo-bar": "s:7:\"foo-bar\";",
	}

	for value, expectation := range trials {
		result := Serialize(value)
		if result != expectation {
			t.Errorf("expected %s, got %s", expectation, result)
		}
	}
}
