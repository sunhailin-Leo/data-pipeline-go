package commiter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMessageCommit_UnknownClient tests that MessageCommit does not panic when given an unknown client type
func TestMessageCommit_UnknownClient(t *testing.T) {
	// Pass a string type client which doesn't match any case in the switch
	assert.NotPanics(t, func() {
		MessageCommit("unknown_client", "test_message", "test_config")
	})
}
