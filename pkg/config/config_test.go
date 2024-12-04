package config

import (
	"testing"

	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/stretchr/testify/assert"
)

func TestStreamConfigGetSourceBySourceName(t *testing.T) {
	// Create a StreamConfig object for testing
	streamConfig := &StreamConfig{
		Source: []*SourceConfig{
			{SourceName: "source1"},
			{SourceName: "source2"},
		},
	}

	// Test existing source name
	source1 := streamConfig.GetSourceBySourceName("source1")
	assert.NotNil(t, source1)
	assert.Equal(t, "source1", source1.SourceName)

	// Test non-existent source name
	source3 := streamConfig.GetSourceBySourceName("source3")
	assert.Nil(t, source3)
}

func TestStreamConfigGetSourceTagBySourceName(t *testing.T) {
	// Create a StreamConfig object for testing
	streamConfig := &StreamConfig{
		Source: []*SourceConfig{
			{SourceName: "source1", Type: "type1"},
			{SourceName: "source2", Type: "type2"},
		},
	}

	// Test existing source name
	tag := streamConfig.GetSourceTagBySourceName("source1")
	assert.Equal(t, "type1", tag)

	// Test non-existent source name
	tag = streamConfig.GetSourceTagBySourceName("source3")
	assert.Equal(t, "", tag)
}

func TestStreamConfigGetSinkBySinkName(t *testing.T) {
	// Create a StreamConfig object for testing
	streamConfig := &StreamConfig{
		Sink: []*SinkConfig{
			{
				SinkName: "sink1",
			},
			{
				SinkName: "sink2",
			},
		},
	}

	// Test existing source name
	sink1 := streamConfig.GetSinkBySinkName("sink1")
	assert.NotNil(t, sink1)

	// Test non-existent source name
	sink3 := streamConfig.GetSinkBySinkName("sink3")
	assert.Nil(t, sink3)
}

func TestStreamConfigGetSinkTagBySinkName(t *testing.T) {
	// Create a StreamConfig object for testing
	streamConfig := &StreamConfig{
		Sink: []*SinkConfig{
			{
				Type:     "type1",
				SinkName: "sink1",
			},
			{
				Type:     "type2",
				SinkName: "sink2",
			},
		},
	}

	// Test existing source name
	tag := streamConfig.GetSinkTagBySinkName("sink1")
	assert.Equal(t, "type1", tag)

	// Test non-existent source name
	tag = streamConfig.GetSinkTagBySinkName("sink3")
	assert.Equal(t, "", tag)
}

func TestTransformSchema(t *testing.T) {
	// Test cases
	testCases := []struct {
		TransformSchema
		wantErr bool
	}{
		// IsIgnore is true, SinkKey can be empty
		{
			TransformSchema: TransformSchema{
				SourceKey:  "sourceKey",
				SinkKey:    "",
				Converter:  "toInt",
				IsIgnore:   true,
				SourceName: "sourceName",
				SinkName:   "sinkName",
			},
			wantErr: false,
		},
		// IsIgnore is true, SinkKey can have a value
		{
			TransformSchema: TransformSchema{
				SourceKey:  "sourceKey",
				SinkKey:    "sinkKey",
				Converter:  "toInt",
				IsIgnore:   true,
				SourceName: "sourceName",
				SinkName:   "sinkName",
			},
			wantErr: false,
		},
		// IsIgnore is false, SinkKey cannot be empty
		{
			TransformSchema: TransformSchema{
				SourceKey:  "sourceKey",
				SinkKey:    "",
				Converter:  "toInt",
				IsIgnore:   false,
				SourceName: "sourceName",
				SinkName:   "sinkName",
			},
			wantErr: true,
		},
		// IsIgnore is false, SinkKey having a value is valid
		{
			TransformSchema: TransformSchema{
				SourceKey:  "sourceKey",
				SinkKey:    "sinkKey",
				Converter:  "toInt",
				IsIgnore:   false,
				SourceName: "sourceName",
				SinkName:   "sinkName",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		schema := tc.TransformSchema
		err := vd.Validate(schema)
		assert.Equal(t, tc.wantErr, err)
	}
}
