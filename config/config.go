package config

import (
	log "github.com/sirupsen/logrus"
)

// LoggingConfig defines the configuration for logging
type LoggingConfig struct {
	Level   string
	Format  string
	Outputs []OutputConfig
}

type OutputConfig struct {
	Type       string
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var config LoggingConfig

// InitLogging initializes the logging configuration
func InitLogging(_config LoggingConfig) {
	//pass to here from bussiness module
	config = _config
	// Parse logrus log level
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

	//constrcut hook slice from configuration
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

	// Configure hooks
	log.AddHook(&LevelHook{
		LevelSlices: hookLevels,
	})
}
