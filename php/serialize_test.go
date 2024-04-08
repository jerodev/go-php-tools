package php

import "testing"

type testStruct struct {
	Name     string
	LastName string `php:"last_name"`
	Age      int8
	Traits   []string
}

func testSerialize(t *testing.T, expectation string, data interface{}) {
	valueString, _ := Serialize(data)

	if valueString != expectation {
		t.Errorf("expected %s, got %s", expectation, valueString)
	}
}

func TestSerializeArray(t *testing.T) {
	testSerialize(t, "a:3:{i:0;i:4;i:1;i:5;i:2;i:6;}", [3]int{4, 5, 6})
	testSerialize(t, "a:2:{i:0;i:1;i:1;i:2;}", []int{1, 2})
	testSerialize(t, "a:3:{i:0;a:2:{i:0;i:1;i:1;i:2;}i:1;a:2:{i:0;i:3;i:1;i:4;}i:2;a:1:{i:0;i:5;}}", [][]int{{1, 2}, {3, 4}, {5}})
}

func TestSerializeMap(t *testing.T) {
	// Remember, map keys are sorted alphabetically because otherwise the order cannot be predicted
	testSerialize(t, "a:2:{s:4:\"That\";i:18;s:4:\"This\";i:7;}", map[string]int{"This": 7, "That": 18})
	testSerialize(t, "a:3:{s:5:\"Maybe\";s:9:\"Misschien\";s:2:\"No\";s:3:\"Nee\";s:3:\"Yes\";s:2:\"Ja\";}", map[string]string{"Yes": "Ja", "No": "Nee", "Maybe": "Misschien"})
}

func TestSerializeObject(t *testing.T) {
	obj := testStruct{
		Name:     "Foo",
		LastName: "Bar",
		Age:      38,
		Traits: []string{
			"Fast",
			"Slow",
		},
	}

	testSerialize(t, "O:10:\"testStruct\":4:{s:4:\"Name\";s:3:\"Foo\";s:9:\"last_name\";s:3:\"Bar\";s:3:\"Age\";i:38;s:6:\"Traits\";a:2:{i:0;s:4:\"Fast\";i:1;s:4:\"Slow\";}}", obj)
}

func TestSerializeScalar(t *testing.T) {
	testSerialize(t, "N;", nil)
	testSerialize(t, "b:1;", true)
	testSerialize(t, "b:0;", false)
	testSerialize(t, "d:1.23;", 1.23)
	testSerialize(t, "i:123;", 123)
	testSerialize(t, "s:7:\"foo-bar\";", "foo-bar")
}

func testUnserialize(t *testing.T, data string, expectation interface{}) {
	var destination interface{}
	err := Unserialize(data, &destination)
	if err != nil {
		t.Fatal(err)
	}

	if destination != expectation {
		t.Errorf("Expected %v, got %v", expectation, destination)
	}
}

func TestUnserializeScalar(t *testing.T) {
	testUnserialize(t, "i:3;", 3)
	testUnserialize(t, "d:3.14;", 3.14)
}
