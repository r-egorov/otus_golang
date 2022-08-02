package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level, OutPath string
	// TODO
}

const configType = "toml"

const (
	defaultLogLevel = "INFO"
)

func NewConfig(configFilePath string) Config {
	viper.SetDefault("logger", map[string]string{
		"level": defaultLogLevel,
		"file":  "stdout",
	})
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("can't read config, err: %w", err))
	}

	loggerMap := viper.GetStringMapString("logger")

	logger := LoggerConf{
		Level:   loggerMap["level"],
		OutPath: loggerMap["file"],
	}
	return Config{
		Logger: logger,
	}
}

func (c Config) GetLogWriter() (out *os.File, outClose func() error) {
	var err error

	switch c.Logger.OutPath {
	case "stdout":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	default:
		out, err = os.OpenFile(c.Logger.OutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			panic(fmt.Errorf("fatal: log file %s, err: %w", c.Logger.OutPath, err))
		}
	}
	outClose = func() error { return out.Close() }
	return
}
