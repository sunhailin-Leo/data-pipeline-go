package commiter

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// rocketmqMessageCommit RocketMQ message commit
func rocketmqMessageCommit(client rocketmq.PullConsumer, record *primitive.MessageQueue, configName string, offsets int64) {
	updateErr := client.UpdateOffset(record, offsets)
	if updateErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[RocketMQ-Source][Current config: " + configName + "]RocketMQ consumer record commit error! Error reason: " + updateErr.Error())
	}
}
