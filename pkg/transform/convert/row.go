package convert

import (
	"strings"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/models"
)

// RowMode row mode
func RowMode(beforeConvertData *models.TransformBeforeConvert, transformConfig config.TransformConfig) *models.TransformAfterConvert {
	resultData := make(map[string][]any)
	afterConvertData := &models.TransformAfterConvert{SourceOutput: beforeConvertData.SourceOutput}
	// If a delimiter is configured, then it will be processed after being sliced according to the delimiter first
	if transformConfig.RowSeparator == "" {
		generateConvertResultData(transformConfig.Schemas[0].SinkName, resultData,
			CastTypes(beforeConvertData.BeforeConvertData, transformConfig.Schemas[0].Converter))
		afterConvertData.AfterConvertData = resultData
		return afterConvertData
	}

	for i, value := range strings.Split(beforeConvertData.BeforeConvertData.(string), transformConfig.RowSeparator) {
		generateConvertResultData(transformConfig.Schemas[i].SinkName, resultData,
			CastTypes(value, transformConfig.Schemas[i].Converter))
	}
	afterConvertData.AfterConvertData = resultData
	return afterConvertData
}
