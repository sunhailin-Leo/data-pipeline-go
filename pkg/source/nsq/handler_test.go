package source

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	baseSource "github.com/sunhailin-Leo/data-pipeline-go/pkg/source"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestNSQSourceHandler(t *testing.T) {
	initLogger()

	handler := &NSQSourceHandler{
		BaseSource: baseSource.BaseSource{
			StreamName:      "test-stream",
			SourceAliasName: "nsq-1",
			SourceConfig: &config.SourceConfig{
				Type:       utils.SourceNSQTagName,
				SourceName: "nsq-1",
				NSQ: config.NSQSourceConfig{
					Topic:   "test-topic",
					Channel: "test-channel",
				},
			},
		},
		sourceTopic:   "test-topic",
		sourceChannel: "test-channel",
	}

	assert.Equal(t, utils.SourceNSQTagName, handler.SourceName())
	assert.Equal(t, "test-topic", handler.SourceTopic())
}
