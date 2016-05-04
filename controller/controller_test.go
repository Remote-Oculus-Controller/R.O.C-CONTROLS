package controller

import (
	"testing"
	"reflect"
	"R.O.C-CONTROLS"
)

func TestParsing(t *testing.T) {

	expected := map[string]roc.Cmd{
		"z":{"forward",1, 0, 255, 127},
		"s":{"stop", 2, 0, 0, 0},
		"q":{"head_xm", 6, 0, 0, 127},
		"d":{"head_xp", 5, 0, 0, 127},
	}
	c := Controller{}
	c.mapControl("../config/ds3_map.json")
	eq := reflect.DeepEqual(c.cmap, expected)
	if !eq {
		t.Error("Error compared map are not the same:\n", expected, c.cmap)
	}
}