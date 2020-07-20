package config

import (
	"strings"

	"github.com/quangdangfit/gosdk/utils/logger"

	"github.com/spf13/viper"
)

func LoadConfig(configFile string) {
	logger.Infof("load config from file: %s", configFile)

	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Panicw("failed to read in config", "error", err)
	}
}
