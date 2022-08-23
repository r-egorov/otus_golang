package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HttpServer ServerConf
	GrpcServer ServerConf
}

type LoggerConf struct {
	Level, OutPath string
}

type StorageConf struct {
	StorageType, User, Password, DBName, Host, Port string
}

type ServerConf struct {
	Host, Port string
}

const configType = "toml"

const (
	DefaultLogLevel     = "INFO"
	InmemoryStorageType = "inmemory"
	PSQLStorageType     = "psql"
)

func NewConfig(configFilePath string) (Config, error) {
	viper.SetDefault("logger", map[string]string{
		"level": DefaultLogLevel,
		"file":  "stdout",
	})
	viper.SetDefault("storage", map[string]string{
		"storage_type": InmemoryStorageType,
	})
	viper.SetDefault("http", map[string]string{
		"host": "localhost",
		"port": "8000",
	})
	viper.SetDefault("grpc", map[string]string{
		"host": "localhost",
		"port": "9000",
	})

	viper.SetConfigType(configType)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("can't read config, err: %w", err)
	}

	storage, err := parseStorageMap(viper.GetStringMapString("storage"))
	if err != nil {
		return Config{}, err
	}
	logger := parseLoggerMap(viper.GetStringMapString("logger"))
	httpServer := parseServerMap(viper.GetStringMapString("http"))
	grpcServer := parseServerMap(viper.GetStringMapString("grpc"))
	return Config{
		Logger:     logger,
		Storage:    storage,
		HttpServer: httpServer,
		GrpcServer: grpcServer,
	}, nil
}

func parseLoggerMap(loggerMap map[string]string) LoggerConf {
	return LoggerConf{
		Level:   loggerMap["level"],
		OutPath: loggerMap["file"],
	}
}

func parseStorageMap(storageMap map[string]string) (StorageConf, error) {
	if storageMap["storage_type"] == "inmemory" {
		return StorageConf{StorageType: "inmemory"}, nil
	}
	getFromConfig := func(key string) (string, error) {
		value, ok := storageMap[key]
		if !ok {
			return "", fmt.Errorf("no %s specified for postgres in config", key)
		}
		return value, nil
	}
	user, err := getFromConfig("user")
	if err != nil {
		return StorageConf{}, err
	}
	password, err := getFromConfig("password")
	if err != nil {
		return StorageConf{}, err
	}
	dbName, err := getFromConfig("db")
	if err != nil {
		return StorageConf{}, err
	}
	host, err := getFromConfig("host")
	if err != nil {
		return StorageConf{}, err
	}
	port, err := getFromConfig("port")
	if err != nil {
		return StorageConf{}, err
	}
	return StorageConf{
		StorageType: "psql",
		User:        user,
		Password:    password,
		DBName:      dbName,
		Host:        host,
		Port:        port,
	}, nil
}

func parseServerMap(serverMap map[string]string) ServerConf {
	return ServerConf{
		Host: serverMap["host"],
		Port: serverMap["port"],
	}
}
