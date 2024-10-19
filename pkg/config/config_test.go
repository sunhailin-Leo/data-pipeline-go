package config

import (
	"testing"

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
