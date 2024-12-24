package php

import (
	"testing"
)

type testStruct struct {
	Name     string
	LastName string `php:"last_name"`
	Age      int8
	Traits   []string
}

type testSimpleStruct struct {
	Name string
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
	testSerialize(t, "a:0:{}", map[string]string{})

	// Remember, map keys are sorted alphabetically because otherwise the order cannot be predicted
	testSerialize(t, "a:2:{s:4:\"That\";i:18;s:4:\"This\";i:7;}", map[string]int{"This": 7, "That": 18})
	testSerialize(t, "a:3:{s:5:\"Maybe\";s:9:\"Misschien\";s:2:\"No\";s:3:\"Nee\";s:3:\"Yes\";s:2:\"Ja\";}", map[string]string{"Yes": "Ja", "No": "Nee", "Maybe": "Misschien"})

	// Special key and value types
	testSerialize(t, "a:2:{d:3.5;s:3:\"Bar\";d:8.8;s:3:\"Foo\";}", map[float64]string{8.8: "Foo", 3.5: "Bar"})
	testSerialize(t, "a:2:{b:1;O:16:\"testSimpleStruct\":1:{s:4:\"Name\";s:3:\"Foo\";}b:0;O:16:\"testSimpleStruct\":1:{s:4:\"Name\";s:3:\"Bar\";}}", map[bool]testSimpleStruct{true: {"Foo"}, false: {"Bar"}})
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

func TestWithStructNames(t *testing.T) {
	WithStructNames(map[string]string{
		"testSimpleStruct": `App\Foo`,
	})

	testSerialize(t, "O:7:\"App\\Foo\":1:{s:4:\"Name\";s:7:\"Jerodev\";}", testSimpleStruct{"Jerodev"})
}

func TestSerializeScalar(t *testing.T) {
	testSerialize(t, "N;", nil)
	testSerialize(t, "b:0;", false)
	testSerialize(t, "b:1;", true)
	testSerialize(t, "d:1.23;", 1.23)
	testSerialize(t, "i:123;", 123)
	testSerialize(t, "s:7:\"foo-bar\";", "foo-bar")
}
