package commiter

import (
	"github.com/nsqio/go-nsq"
)

// nsqMessageCommit 确认一条已消费的 NSQ 消息。
func nsqMessageCommit(message *nsq.Message, _ string) {
	if message != nil {
		message.Finish()
	}
}
