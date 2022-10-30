package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/rmq"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configFilePath)
	if err != nil {
		log.Fatalf("config: %s", err.Error())
	}

	// Configure logger
	logOut, logOutClose := getLogWriter(conf)
	defer func() {
		if err := logOutClose(); err != nil {
			panic(err)
		}
	}()
	logg := logger.New(logOut, conf.Logger.Level)

	// Configure global context to watch the signals
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rabbit := rmq.New(conf.AMQP.URI, conf.AMQP.Queue)
	err = rabbit.Connect()
	if err != nil {
		logg.Fatal("failed to connect to RabbitMQ")
	}
	defer func() {
		err := rabbit.Close()
		if err != nil {
			logg.Fatal("failed to close connection RabbitMQ")
		}
	}()

	msgs, err := rabbit.Consume()
	if err != nil {
		logg.Fatal("failed to consume amqp")
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var notification storage.Notification
				if err := json.Unmarshal(msg.Body, &notification); err != nil {
					logg.Error(fmt.Sprintf("failed to unmarshal message: %v", err))
				}
				logg.Info(notification.String())
			}
		}
	}(ctx)

	logg.Info("started sender...")
	<-ctx.Done()
}

func getLogWriter(c config.Config) (out *os.File, outClose func() error) {
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
