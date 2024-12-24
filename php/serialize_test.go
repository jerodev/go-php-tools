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

func testSerialize(t *testing.T, expectation string, data interface{}) error {
	valueString, err := Serialize(data)

	if valueString != expectation {
		t.Errorf("expected %s, got %s", expectation, valueString)
	}

	return err
}

func TestSerializeArray(t *testing.T) {
	four, five := 4, 5

	testSerialize(t, "a:0:{}", []string{})
	testSerialize(t, "a:3:{i:0;i:4;i:1;i:5;i:2;i:6;}", [3]int{4, 5, 6})
	testSerialize(t, "a:2:{i:0;i:1;i:1;i:2;}", []int{1, 2})
	testSerialize(t, "a:2:{i:0;i:4;i:1;i:5;}", []*int{&four, &five})
	testSerialize(t, "a:3:{i:0;a:2:{i:0;i:1;i:1;i:2;}i:1;a:2:{i:0;i:3;i:1;i:4;}i:2;a:1:{i:0;i:5;}}", [][]int{{1, 2}, {3, 4}, {5}})
}

func TestSerializeMap(t *testing.T) {
	testSerialize(t, "a:0:{}", map[string]string{})

	// Remember, map keys are sorted alphabetically because otherwise the order cannot be predicted
	testSerialize(t, "a:2:{s:4:\"That\";i:18;s:4:\"This\";i:7;}", map[string]int{"This": 7, "That": 18})
	testSerialize(t, "a:3:{s:5:\"Maybe\";s:9:\"Misschien\";s:2:\"No\";s:3:\"Nee\";s:3:\"Yes\";s:2:\"Ja\";}", map[string]string{"Yes": "Ja", "No": "Nee", "Maybe": "Misschien"})

	// Special key and value types
	testSerialize(t, "a:2:{i:3;s:3:\"Bar\";i:8;s:3:\"Foo\";}", map[float64]string{8.8: "Foo", 3.5: "Bar"})
	testSerialize(t, "a:2:{i:0;O:16:\"testSimpleStruct\":1:{s:4:\"Name\";s:3:\"Bar\";}i:1;O:16:\"testSimpleStruct\":1:{s:4:\"Name\";s:3:\"Foo\";}}", map[bool]testSimpleStruct{true: {"Foo"}, false: {"Bar"}})

	// Pointer key/values
	seven := 7
	testSerialize(t, "a:1:{i:7;i:77;}", map[*int]int{&seven: 77})
	testSerialize(t, "a:1:{i:7;i:7;}", map[*int]*int{&seven: &seven})
	testSerialize(t, "a:1:{s:0:\"\";i:77;}", map[*int]int{nil: 77})
}

func TestInvalidArrayKeys(t *testing.T) {
	err := testSerialize(t, "", map[testSimpleStruct]string{
		testSimpleStruct{"Foo"}: "bar",
	})

	if err == nil {
		t.Error("struct array keys should not be allowed but are")
	}
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
