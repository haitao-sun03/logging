package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LevelHook implements logrus.Hook
type LevelHook struct {
	LevelSlices []log.Level
}

// Levels returns the log levels that this hook should be fired for
func (hook *LevelHook) Levels() []log.Level {
	return hook.LevelSlices
}

// Fire is called when a log event occurs
func (hook *LevelHook) Fire(entry *log.Entry) error {
	writer, err := createLogWriterForLevel(entry.Level)
	if err != nil {
		return err
	}
	multiWriter := io.MultiWriter(writer, os.Stdout)
	entry.Logger.Out = multiWriter
	return nil
}

// createLogWriterForLevel creates a log writer for the given log level
func createLogWriterForLevel(level log.Level) (io.Writer, error) {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Find the corresponding output configuration
	var outputConfig *OutputConfig
	for _, output := range config.Outputs {
		if output.Type == "file" && strings.EqualFold(output.Level, level.String()) {
			outputConfig = &output
			break
		}
	}
	if outputConfig == nil {
		return nil, fmt.Errorf("no configuration found for level %s", level.String())
	}

	now := time.Now()
	filename := filepath.Join(logDir, now.Format(outputConfig.Filename))

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    outputConfig.MaxSize,
		MaxBackups: outputConfig.MaxBackups,
		MaxAge:     outputConfig.MaxAge,
		Compress:   outputConfig.Compress,
	}, nil
}
