package convert

import (
	"reflect"
	"testing"

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
				"Clickhouse-1": []any{123},
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
