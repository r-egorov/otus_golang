package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer ServerConf
	GRPCServer ServerConf
	AMQP       AMQPConf
	Scheduler  SchedulerConf
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

type AMQPConf struct {
	URI, Queue string
}

type SchedulerConf struct {
	NotifyBefore, Period time.Duration
}

const configType = "toml"

const (
	DefaultLogLevel     = "INFO"
	DefaultLogOutFile   = "stdout"
	InmemoryStorageType = "inmemory"
	PSQLStorageType     = "psql"
	DefaultHTTPHost     = "localhost"
	DefaultHTTPPort     = "8000"
	DefaultGRPCHost     = "localhost"
	DefaultGRPCPort     = "9000"
	DefaultAMQPUri      = "amqp://guest:guest@localhost:5672"
	DefaultAMQPQueue    = "calendar"
	DefaultRemindBefore = "1h"
	DefaultPeriod       = "10m"
)

func NewConfig(configFilePath string) (Config, error) {
	viper.SetDefault("logger", map[string]string{
		"level": DefaultLogLevel,
		"file":  DefaultLogOutFile,
	})
	viper.SetDefault("storage", map[string]string{
		"storage_type": InmemoryStorageType,
	})
	viper.SetDefault("http", map[string]string{
		"host": DefaultHTTPHost,
		"port": DefaultHTTPPort,
	})
	viper.SetDefault("grpc", map[string]string{
		"host": DefaultGRPCHost,
		"port": DefaultGRPCPort,
	})
	viper.SetDefault("amqp", map[string]string{
		"uri":   DefaultAMQPUri,
		"queue": DefaultAMQPQueue,
	})
	viper.SetDefault("scheduler", map[string]string{
		"notify_before": DefaultRemindBefore,
		"period":        DefaultPeriod,
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
	amqp := parseAMQPMap(viper.GetStringMapString("amqp"))
	scheduler, err := parseSchedulerMap(viper.GetStringMapString("scheduler"))
	if err != nil {
		return Config{}, err
	}
	return Config{
		Logger:     logger,
		Storage:    storage,
		HTTPServer: httpServer,
		GRPCServer: grpcServer,
		AMQP:       amqp,
		Scheduler:  scheduler,
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

func parseAMQPMap(amqpMap map[string]string) AMQPConf {
	return AMQPConf{
		URI:   amqpMap["uri"],
		Queue: amqpMap["queue"],
	}
}

func parseSchedulerMap(schedulerMap map[string]string) (SchedulerConf, error) {
	remindBefore, err := time.ParseDuration(schedulerMap["notify_before"])
	if err != nil {
		return SchedulerConf{}, err
	}
	period, err := time.ParseDuration(schedulerMap["period"])
	if err != nil {
		return SchedulerConf{}, err
	}
	return SchedulerConf{
		NotifyBefore: remindBefore,
		Period:       period,
	}, nil
}
