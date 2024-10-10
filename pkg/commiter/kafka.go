package commiter

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

// kafkaMessageCommit Kafka message commit
func kafkaMessageCommit(client *kgo.Client, record *kgo.Record, configName string) {
	commitErr := client.CommitRecords(context.Background(), record)
	if commitErr != nil {
		logger.Logger.Error(utils.LogServiceName +
			"[Kafka-Source][Current config: " + configName + "]Kafka consumer record commit error! Error reason: " + commitErr.Error())
	}
}
