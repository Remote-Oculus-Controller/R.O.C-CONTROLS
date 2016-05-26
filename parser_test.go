package roc

import (
	"testing"
	"reflect"
)

func TestDecode(t *testing.T)  {

	expected := map[string]interface{}{
		"z":"forward",
		"s":"stop",
		"q":"head_xm",
		"d":"head_xp",
	}
	m, err := DecodeJsonFile("./config/parser_test_decode.json")
	if err != nil {
		t.Error(err)
	}
	eq := reflect.DeepEqual(expected, m)
	if !eq {
		t.Error("Error compared map are not the same:\n", expected, m)
	}
}
