package commiter

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/wagslane/go-rabbitmq"
)

// MessageCommit unified message commit function
func MessageCommit(client, message interface{}, configName string, params ...interface{}) {
	switch c := client.(type) {
	case *kgo.Client:
		kafkaMessageCommit(c, message.(*kgo.Record), configName)
	case rocketmq.PullConsumer:
		rocketmqMessageCommit(c, message.(*primitive.MessageQueue), configName, params[0].(int64))
	case rabbitmq.Delivery:
		rabbitmqMessageCommit(c, configName)
	case pulsar.Consumer:
		pulsarMessageCommit(c, message.(pulsar.Message), configName)
	}
}
