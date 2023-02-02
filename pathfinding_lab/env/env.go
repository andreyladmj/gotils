package env

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ReadOS() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("LOG_LEVEL", "warning")

	logLevel, err := log.ParseLevel(v.GetString("LOG_LEVEL"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	return &Config{
		LogLevel: logLevel,
	}, nil
}
