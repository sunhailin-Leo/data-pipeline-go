package utils

import (
	"reflect"
	"testing"
)

func TestStringToByte(t *testing.T) {
	t.Helper()

	actual := StringToByte("test")
	expected := []byte{116, 101, 115, 116}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("StringToByte failed. Expected: %v, actual: %v", expected, actual)
	}
}

func TestByteToString(t *testing.T) {
	t.Helper()

	actual := ByteToString([]byte{116, 101, 115, 116})
	expected := "test"
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("ByteToString failed. Expected: %v, actual: %v", expected, actual)
	}
}
