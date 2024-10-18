package config

import (
	"os"
	"testing"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/logger"
	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func initLogger() {
	logger.NewZapLogger()
}

func TestLoadConfigFromLocal(t *testing.T) {
	initLogger()

	testConfigFile := "config.json"

	_ = os.Setenv(utils.ConfigFromSourceName, "local")
	_ = os.Setenv(utils.ConfigFromLocalPathEnvName, testConfigFile)
	NewTunnelConfigLoader()

	// TODO Create json and load

}
