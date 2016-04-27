package main

import (
	"testing"
)

func TestGetAngle(t *testing.T) {
	expectedInt := 54
	result := getAngle(75, 500)
	if result != expectedInt {
		t.Fatalf("Expected %s, got %s", expectedInt, result)
	}

	expectedInt = 360
	result = getAngle(500, 500)
	if result != expectedInt {
		t.Fatalf("Expected %s, got %s", expectedInt, result)
	}
}

func TestDiffIsCorrect(t *testing.T) {
	expectedBool := true
	result := diffIsCorrect(750, 753)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}

	expectedBool = false
	result = diffIsCorrect(750, 0)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}

	expectedBool = false
	result = diffIsCorrect(750, -10)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}

	expectedBool = false
	result = diffIsCorrect(750, 1600)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}

	expectedBool = true
	result = diffIsCorrect(750, 750)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}

	expectedBool = false
	result = diffIsCorrect(750, 540)
	if result != expectedBool {
		t.Fatalf("Expected %s, got %s", expectedBool, result)
	}
}