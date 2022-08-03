package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	// TODO
}

type LoggerConf struct {
	Level, OutPath string
}

type StorageConf struct {
	StorageType, User, Password, DBName, Host, Port string
}

const configType = "toml"

const (
	defaultLogLevel     = "INFO"
	inmemoryStorageType = "inmemory"
	psqlStorageType     = "psql"
)

func NewConfig(configFilePath string) Config {
	viper.SetDefault("logger", map[string]string{
		"level": defaultLogLevel,
		"file":  "stdout",
	})
	viper.SetDefault("storage", map[string]string{
		"storage_type": inmemoryStorageType,
	})

	viper.SetConfigType(configType)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("can't read config, err: %w", err))
	}

	storage := parseStorageMap(viper.GetStringMapString("storage"))
	logger := parseLoggerMap(viper.GetStringMapString("logger"))
	return Config{
		Logger:  logger,
		Storage: storage,
	}
}

func parseLoggerMap(loggerMap map[string]string) LoggerConf {
	return LoggerConf{
		Level:   loggerMap["level"],
		OutPath: loggerMap["file"],
	}
}

func parseStorageMap(storageMap map[string]string) StorageConf {
	if storageMap["storage_type"] == "inmemory" {
		return StorageConf{StorageType: "inmemory"}
	}
	getFromConfig := func(key string) string {
		value, ok := storageMap[key]
		if !ok {
			panic(fmt.Errorf("no %s specified for postgres in config", key))
		}
		return value
	}
	user := getFromConfig("user")
	password := getFromConfig("password")
	dbName := getFromConfig("db")
	host := getFromConfig("host")
	port := getFromConfig("port")
	return StorageConf{
		StorageType: "psql",
		User:        user,
		Password:    password,
		DBName:      dbName,
		Host:        host,
		Port:        port,
	}
}
