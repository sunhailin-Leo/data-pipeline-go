package convert

import (
	"reflect"
	"testing"

	"github.com/cloudwego/gjson"
	"github.com/stretchr/testify/assert"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestGenerateConvertResultData(t *testing.T) {
	t.Helper()

	type testCaseUnit struct {
		SinkNames     string
		AfterCaseData any
		Expected      map[string][]any
	}

	testUnits := []testCaseUnit{
		{
			SinkNames:     "Clickhouse-1",
			AfterCaseData: 123,
			Expected: map[string][]any{
				"Clickhouse-1": {123},
			},
		},
		{
			SinkNames:     "Clickhouse-1,Clickhouse-2",
			AfterCaseData: 123,
			Expected: map[string][]any{
				"Clickhouse-1": {123},
				"Clickhouse-2": {123},
			},
		},
	}

	for _, unit := range testUnits {
		actual := make(map[string][]any)
		generateConvertResultData(unit.SinkNames, actual, unit.AfterCaseData)
		if !reflect.DeepEqual(actual, unit.Expected) {
			t.Fatalf("TestGenerateConvertResultData failed. Expected: %v, actual: %v", unit.Expected, actual)
		}
	}
}

func TestJsonPathToMap(t *testing.T) {
	initLogger()

	// normal
	data := []byte(`{"key1": "value1", "key2": "value2"}`)
	paths := []config.TransformJsonPath{
		{SrcField: "", Path: "key1", DestField: "destKey1"},
		{SrcField: "", Path: "key2", DestField: "destKey2"},
	}

	result := JsonPathToMap(data, paths)
	assert.Equal(t, map[string]any{"destKey1": "value1", "destKey2": "value2"}, result)

	// path error
	data = []byte(`{"key1": "value1", "key2": "value2"}`)
	paths = []config.TransformJsonPath{
		{SrcField: "", Path: "invalidPath", DestField: "destKey"},
	}

	result = JsonPathToMap(data, paths)
	assert.Equal(t, map[string]any{}, result)

	// data error
	data = []byte(`invalidData`)
	paths = []config.TransformJsonPath{
		{SrcField: "", Path: "key1", DestField: "destKey1"},
	}

	result = JsonPathToMap(data, paths)
	assert.Equal(t, map[string]any{}, result)
}

func TestJsonToMap(t *testing.T) {
	t.Helper()

	type testCaseUnit struct {
		Data     []byte
		Expected map[string]any
	}

	testUnits := []testCaseUnit{
		{
			Data: []byte(`{"name": "test"}`),
			Expected: map[string]any{
				"name": "test",
			},
		},
	}

	for _, unit := range testUnits {
		actual := JsonToMap(unit.Data)
		if !reflect.DeepEqual(actual, unit.Expected) {
			t.Fatalf("TestJsonToMap failed. Expected: %v, actual: %v", unit.Expected, actual)
		}
	}
}

func TestMapToJson(t *testing.T) {
	t.Helper()

	type testCaseUnit struct {
		Data     map[string]any
		Expected []byte
	}

	testUnits := []testCaseUnit{
		{
			Data: map[string]any{
				"name": "test",
			},
			Expected: []byte(`{"name":"test"}`),
		},
	}

	for _, unit := range testUnits {
		actual := MapToJson(unit.Data)
		if !reflect.DeepEqual(actual, unit.Expected) {
			t.Fatalf("TestMapToJson failed. Expected: %v, actual: %v", unit.Expected, actual)
		}
	}
}

func TestCastTypes(t *testing.T) {
	t.Helper()

	initLogger()

	type testCaseUnit struct {
		Data          any
		ConvertorName string
		Expected      any
	}

	testUnits := []testCaseUnit{
		{
			Data:          "123",
			ConvertorName: "toInt64",
			Expected:      int64(123),
		},
		{
			Data:          "123",
			ConvertorName: "toFloat64",
			Expected:      float64(123),
		},
		{
			Data:          "123",
			ConvertorName: "toFloat32",
			Expected:      float32(123),
		},
		{
			Data:          "true",
			ConvertorName: "toBool",
			Expected:      true,
		},
		{
			Data:          "456",
			ConvertorName: "toInt32",
			Expected:      int32(456),
		},
		{
			Data:          "789",
			ConvertorName: "toInt16",
			Expected:      int16(789),
		},
		{
			Data:          "123",
			ConvertorName: "toInt8",
			Expected:      int8(123),
		},
		{
			Data:          "999",
			ConvertorName: "toInt",
			Expected:      999,
		},
		{
			Data:          "1000",
			ConvertorName: "toUint",
			Expected:      uint(1000),
		},
		{
			Data:          "2000",
			ConvertorName: "toUint64",
			Expected:      uint64(2000),
		},
		{
			Data:          "3000",
			ConvertorName: "toUint32",
			Expected:      uint32(3000),
		},
		{
			Data:          "4000",
			ConvertorName: "toUint16",
			Expected:      uint16(4000),
		},
		{
			Data:          "200",
			ConvertorName: "toUint8",
			Expected:      uint8(200),
		},
		{
			Data:          "hello",
			ConvertorName: "toString",
			Expected:      "hello",
		},
		{
			Data:          123,
			ConvertorName: "toString",
			Expected:      "123",
		},
		{
			Data:          "hello",
			ConvertorName: "",
			Expected:      "hello",
		},
		{
			Data:          "invalid",
			ConvertorName: "toInt",
			Expected:      nil,
		},
		{
			Data:          "test",
			ConvertorName: "unknown",
			Expected:      nil,
		},
	}

	for _, unit := range testUnits {
		actual := CastTypes(unit.Data, unit.ConvertorName)
		if !reflect.DeepEqual(actual, unit.Expected) {
			t.Fatalf("TestCastTypes failed. ConvertorName: %s, Expected: %v, actual: %v", unit.ConvertorName, unit.Expected, actual)
		}
	}
}

func TestCastTypesDefaultValue(t *testing.T) {
	// empty string
	result := CastTypesDefaultValue("")
	assert.Nil(t, result)

	// test "toBool"
	result = CastTypesDefaultValue("toBool")
	assert.Equal(t, false, result)

	// test "toFloat64" and "toFloat32"
	result = CastTypesDefaultValue("toFloat64")
	assert.Equal(t, 0.0, result)
	result = CastTypesDefaultValue("toFloat32")
	assert.Equal(t, 0.0, result)

	// test other int types
	result = CastTypesDefaultValue("toInt64")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toInt32")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toInt16")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toInt8")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toInt")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toUint")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toUint64")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toUint32")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toUint16")
	assert.Equal(t, 0, result)
	result = CastTypesDefaultValue("toUint8")
	assert.Equal(t, 0, result)

	// test "toString"
	result = CastTypesDefaultValue("toString")
	assert.Equal(t, "", result)

	// test other
	result = CastTypesDefaultValue("other")
	assert.Nil(t, result)
}

func TestCastFunctionNameToFunctionResult(t *testing.T) {
	// Test scenario 1: The input is "$.UUID()", and the expected return is the simulated UUID
	result1 := CastFunctionNameToFunctionResult("$.UUID()")
	assert.NotEqual(t, nil, result1)

	// Test scenario 2: The input is something else, and the expected return is the input value itself
	result2 := CastFunctionNameToFunctionResult("other-value")
	assert.Equal(t, "other-value", result2)
}

func NewJsonPathToMap(data []byte, paths []config.TransformJsonPath) map[string]any {
	result := make(map[string]any)
	dataStr := gjson.ParseBytes(data).String()
	for _, path := range paths {
		pathRes := gjson.Get(dataStr, path.Path).Value()
		if pathRes == nil {
			continue
		}
		result[path.DestField] = pathRes
	}

	return result
}

func TestJsonPathQueriesToMap2(t *testing.T) {

	data := []byte(`{
  "store": {
    "book": [
      {
        "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      {
        "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honor",
        "price": 12.99
      },
      {
        "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      {
        "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`)
	paths := []config.TransformJsonPath{
		{SrcField: "", Path: "store.book.#.author", DestField: "author"},
		{SrcField: "", Path: "store.book.#(price>15)#.price", DestField: "price"},
	}
	result := NewJsonPathToMap(data, paths)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville", "J. R. R. Tolkien"}, result["author"])
	assert.Equal(t, []interface{}{22.99}, result["price"])
}

// Benchmark functions
func BenchmarkCastTypes(b *testing.B) {
	initLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CastTypes("123", "toInt64")
	}
}

func BenchmarkJsonToMap(b *testing.B) {
	b.ReportAllocs()
	data := []byte(`{"name": "test", "age": 25, "active": true}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JsonToMap(data)
	}
}

func BenchmarkMapToJson(b *testing.B) {
	b.ReportAllocs()
	data := map[string]any{"name": "test", "age": 25, "active": true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToJson(data)
	}
}

func BenchmarkCastTypesDefaultValue(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CastTypesDefaultValue("toInt64")
	}
}

func BenchmarkJsonPathToMap(b *testing.B) {
	initLogger()
	b.ReportAllocs()
	data := []byte(`{"key1": "value1", "key2": "value2", "key3": "value3"}`)
	paths := []config.TransformJsonPath{
		{SrcField: "", Path: "key1", DestField: "destKey1"},
		{SrcField: "", Path: "key2", DestField: "destKey2"},
		{SrcField: "", Path: "key3", DestField: "destKey3"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		JsonPathToMap(data, paths)
	}
}

func TestJsonToMap_InvalidJson(t *testing.T) {
	initLogger()
	data := []byte(`invalid json`)
	result := JsonToMap(data)
	assert.Nil(t, result)
}

func TestMapToJson_Success(t *testing.T) {
	data := map[string]any{"name": "test", "age": 25}
	result := MapToJson(data)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
}

func TestMapToJson_Error(t *testing.T) {
	// Test with function value which cannot be marshaled by sonic
	data := map[string]any{"func": func() {}}
	result := MapToJson(data)
	assert.Nil(t, result)
}
