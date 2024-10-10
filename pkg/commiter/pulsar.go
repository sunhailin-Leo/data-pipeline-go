package commiter

import (
	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// pulsarMessageCommit Pulsar message commit
func pulsarMessageCommit(client pulsar.Consumer, record pulsar.Message, configName string) {
	ackErr := client.Ack(record)
	if ackErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[Pulsar-Source][Current config: " + configName + "]Pulsar consumer record commit error! Error reason: " + ackErr.Error())
	}
}
