package config

import (
	log "github.com/sirupsen/logrus"
)

// LoggingConfig defines the configuration for logging
type LoggingConfig struct {
	Level   string         `mapstructure:"level"`
	Format  string         `mapstructure:"format"`
	Outputs []OutputConfig `mapstructure:"output"`
}

type OutputConfig struct {
	Type       string `mapstructure:"type"`
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

var config LoggingConfig

// InitLogging initializes the logging configuration
func InitLogging(_config LoggingConfig) {
	// Set viper to read config file
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath("./config")
	// if err := viper.ReadInConfig(); err != nil {
	// 	panic(fmt.Errorf("fatal error config file: %s", err))
	// }

	// if err := viper.UnmarshalKey("logging", &config); err != nil {
	// 	panic(fmt.Errorf("unable to decode into struct, %v", err))
	// }

	config = _config
	// Parse log level
	level, err := log.ParseLevel(config.Level)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	// Configure hooks
	hookLevels := []log.Level{}
	for _, output := range config.Outputs {
		if output.Type == "file" {
			level, err := log.ParseLevel(output.Level)
			if err != nil {
				panic(err)
			}
			hookLevels = append(hookLevels, level)
		}
	}

	log.AddHook(&LevelHook{
		LevelSlices: hookLevels,
	})
}
