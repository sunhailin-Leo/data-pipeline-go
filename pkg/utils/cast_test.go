package utils

import (
	"reflect"
	"testing"
)

func TestStringToByte(t *testing.T) {
	t.Helper()

	t.Run("normal string", func(t *testing.T) {
		actual := StringToByte("test")
		expected := []byte{116, 101, 115, 116}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("StringToByte failed. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		actual := StringToByte("")
		if len(actual) != 0 {
			t.Fatalf("StringToByte failed for empty string. Expected length 0, actual: %v", actual)
		}
	})

	t.Run("string with special characters", func(t *testing.T) {
		actual := StringToByte("hello world!")
		expected := []byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("StringToByte failed for special characters. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("string with unicode", func(t *testing.T) {
		actual := StringToByte("你好")
		expected := []byte{228, 189, 160, 229, 165, 189}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("StringToByte failed for unicode. Expected: %v, actual: %v", expected, actual)
		}
	})
}

func TestByteToString(t *testing.T) {
	t.Helper()

	t.Run("normal bytes", func(t *testing.T) {
		actual := ByteToString([]byte{116, 101, 115, 116})
		expected := "test"
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ByteToString failed. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("empty bytes", func(t *testing.T) {
		actual := ByteToString([]byte{})
		expected := ""
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ByteToString failed for empty bytes. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("bytes with special characters", func(t *testing.T) {
		actual := ByteToString([]byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33})
		expected := "hello world!"
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ByteToString failed for special characters. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("bytes with unicode", func(t *testing.T) {
		actual := ByteToString([]byte{228, 189, 160, 229, 165, 189})
		expected := "你好"
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("ByteToString failed for unicode. Expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("roundtrip", func(t *testing.T) {
		original := "test string with unicode: 你好世界"
		converted := ByteToString(StringToByte(original))
		if !reflect.DeepEqual(converted, original) {
			t.Fatalf("Roundtrip failed. Expected: %v, actual: %v", original, converted)
		}
	})
}

// Benchmark functions
func BenchmarkStringToByte(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToByte("test string for benchmark")
	}
}

func BenchmarkByteToString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	data := []byte("test string for benchmark")
	for i := 0; i < b.N; i++ {
		ByteToString(data)
	}
}
