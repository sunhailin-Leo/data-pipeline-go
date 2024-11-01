package convert

import (
	"reflect"
	"testing"

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
