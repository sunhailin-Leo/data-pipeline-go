package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func TestGetEncodeConfig(t *testing.T) {
	configWithServiceName := getEncodeConfig(true)
	assert.Equal(t, "logger", configWithServiceName.NameKey)
	assert.Equal(t, "linenum", configWithServiceName.CallerKey)
	assert.Equal(t, "level", configWithServiceName.LevelKey)

	configWithoutServiceName := getEncodeConfig(false)
	assert.Equal(t, "logger", configWithoutServiceName.NameKey)
	assert.Equal(t, "linenum", configWithoutServiceName.CallerKey)
	assert.Equal(t, "level", configWithoutServiceName.LevelKey)
}

func TestGetRollingConfig(t *testing.T) {
	filename := "test.log"
	expectedLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    utils.LogMaxSize,
		MaxAge:     utils.LogMaxAge,
		MaxBackups: utils.LogMaxBackups,
		Compress:   true,
	}
	actualLogger := getRollingConfig(filename)
	assert.Equal(t, expectedLogger, actualLogger)
}
