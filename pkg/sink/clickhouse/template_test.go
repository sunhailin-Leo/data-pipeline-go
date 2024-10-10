package sink

import (
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
)

func TestRenderClickhouseCreateTableTemplate(t *testing.T) {
	sinkConfig := config.ClickhouseSinkConfig{
		Database:  "ai_group",
		TableName: "test",
		Columns: []config.ClickhouseTableColumn{
			{Name: "col1", Type: "Int32", Comment: "col1 test"},
			{Name: "col2", Type: "Float32"},
			{Name: "col3", Type: "String"},
		},
		Engine:     "MergeTree()",
		Partition:  []string{"col1"},
		PrimaryKey: []string{"col1"},
		OrderBy:    []string{"col1"},
		Comment:    "test",
		Settings:   []string{"storage_policy = 'policy_data8'"},
	}

	template, err := RenderClickhouseCreateTableTemplate("autoCreateTable", sinkConfig)
	if err != nil {
		panic(err)
	}
	println(template)
}
