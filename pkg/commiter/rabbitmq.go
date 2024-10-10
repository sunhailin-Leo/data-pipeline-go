package commiter

import (
	"github.com/wagslane/go-rabbitmq"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// rabbitmqMessageCommit RabbitMQ message commit
func rabbitmqMessageCommit(client rabbitmq.Delivery, configName string) {
	ackErr := client.Ack(true)
	if ackErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[RabbitMQ-Source][Current config: " + configName + "]RabbitMQ consumer record commit error! Error reason: " + ackErr.Error())
	}
}
