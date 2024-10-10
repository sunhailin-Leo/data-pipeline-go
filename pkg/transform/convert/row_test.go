package convert

import (
	"reflect"
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
)

func TestRowMode(t *testing.T) {
	t.Helper()

	initLogger()

	type testCaseUnit struct {
		InputData *models.TransformBeforeConvert
		Config    config.TransformConfig
		Expected  *models.TransformAfterConvert
	}

	testUnits := []testCaseUnit{
		{
			InputData: &models.TransformBeforeConvert{
				SourceOutput:      &models.SourceOutput{},
				BeforeConvertData: "123",
			},
			Config: config.TransformConfig{
				Mode: "row",
				Schemas: []config.TransformSchema{
					{
						Converter:  "toInt",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
				},
			},
			Expected: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {123},
					"Clickhouse-2": {123},
				},
			},
		},
		{
			InputData: &models.TransformBeforeConvert{
				SourceOutput:      &models.SourceOutput{},
				BeforeConvertData: "123,abc",
			},
			Config: config.TransformConfig{
				Mode: "row",
				Schemas: []config.TransformSchema{
					{
						Converter:  "toInt",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
					{
						Converter:  "toString",
						SourceName: "Kafka-1",
						SinkName:   "Clickhouse-1,Clickhouse-2",
					},
				},
				RowSeparator: ",",
			},
			Expected: &models.TransformAfterConvert{
				SourceOutput: &models.SourceOutput{},
				AfterConvertData: map[string][]any{
					"Clickhouse-1": {123, "abc"},
					"Clickhouse-2": {123, "abc"},
				},
			},
		},
	}

	for _, unit := range testUnits {
		actual := RowMode(unit.InputData, unit.Config)
		if !reflect.DeepEqual(actual.AfterConvertData, unit.Expected.AfterConvertData) {
			t.Fatalf("TestRowMode failed. Expected: %v, actual: %v", unit.Expected.AfterConvertData, actual.AfterConvertData)
		}
	}
}
