package inquiry_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stevenmatthewt/inquiry"
)

func TestSimpleStrings(t *testing.T) {
	type output struct {
		Field1 string `query:"field1"`
		Field2 string `query:"field2"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"val1.1",
		},
		"field2": []string{
			"val2.1",
		},
	}

	var out output
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err != nil {
		t.Error(err)
	}

	if out.Field1 != "val1.1" {
		t.Errorf("Field1 incorrect value: %s", out.Field1)
	}

	if out.Field2 != "val2.1" {
		t.Errorf("Field2 incorrect value: %s", out.Field2)
	}
}

func TestSimpleInts(t *testing.T) {
	type output struct {
		Field1 int `query:"field1"`
		Field2 int `query:"field2"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"23",
		},
		"field2": []string{
			"-42",
		},
	}

	var out output
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err != nil {
		t.Error(err)
	}

	if out.Field1 != 23 {
		t.Errorf("Field1 incorrect value: %d", out.Field1)
	}

	if out.Field2 != -42 {
		t.Errorf("Field2 incorrect value: %d", out.Field2)
	}
}

func TestSimpleIntAndString(t *testing.T) {
	type output struct {
		Field1 int    `query:"field1"`
		Field2 string `query:"field2"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"23",
		},
		"field2": []string{
			"val2.1",
		},
	}

	var out output
	fmt.Printf("correct address %p\n", &out.Field1)
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err != nil {
		t.Error(err)
	}

	if out.Field1 != 23 {
		t.Errorf("Field1 incorrect value: %d", out.Field1)
	}

	if out.Field2 != "val2.1" {
		t.Errorf("Field2 incorrect value: %s", out.Field2)
	}
}

func TestOverflow(t *testing.T) {
	type output struct {
		Field1 int8 `query:"field1"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"1000",
		},
	}

	var out output
	fmt.Printf("correct address %p\n", &out.Field1)
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err == nil {
		t.Fatal("non-nil error returned")
	}

	if !strings.Contains(err.Error(), "overflow") {
		t.Errorf("error of wrong type: %s", err.Error())
	}
}

func TestUnderflow(t *testing.T) {
	type output struct {
		Field1 uint64 `query:"field1"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"-1",
		},
	}

	var out output
	fmt.Printf("correct address %p\n", &out.Field1)
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err == nil {
		t.Fatal("non-nil error returned")
	}

	if !strings.Contains(err.Error(), "underflow") {
		t.Errorf("error of wrong type: %s", err.Error())
	}
}

func TestSimpleFloats(t *testing.T) {
	type output struct {
		Field1 float32 `query:"field1"`
		Field2 float64 `query:"field2"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"23.0",
		},
		"field2": []string{
			"-42.4567",
		},
	}

	var out output
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err != nil {
		t.Error(err)
	}

	if out.Field1 != 23.0 {
		t.Errorf("Field1 incorrect value: %d", out.Field1)
	}

	if out.Field2 != -42.4567 {
		t.Errorf("Field2 incorrect value: %d", out.Field2)
	}
}

func TestSimpleSlice(t *testing.T) {
	type output struct {
		Field1 []int    `query:"field1"`
		Field2 []string `query:"field2"`
	}
	qsMap := map[string][]string{
		"field1": []string{
			"12",
			"34",
			"5",
			"6",
			"789",
		},
		"field2": []string{
			"-42.0",
			"hello",
		},
	}

	var out output
	err := inquiry.UnmarshalMap(qsMap, &out)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(out.Field1, []int{12, 34, 5, 6, 789}) {
		t.Errorf("Field1 incorrect value: %d", out.Field1)
	}

	if !reflect.DeepEqual(out.Field2, []string{"-42.0", "hello"}) {
		t.Errorf("Field2 incorrect value: %d", out.Field2)
	}
}
