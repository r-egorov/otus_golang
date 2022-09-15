package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	conn, err := amqp.Dial(conf.AMQP.URI)
	if err != nil {
		logg.Fatal("failed to connect to RabbitMQ")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logg.Error(fmt.Sprintf("can't close RabbitMQ conn: %v", err))
		}
	}()

	amqpCh, err := conn.Channel()
	if err != nil {
		logg.Fatal("failed to open RabbitMQ channel")
	}
	defer func() {
		err := amqpCh.Close()
		if err != nil {
			logg.Error(fmt.Sprintf("can't close RabbitMQ chan: %v", err))
		}
	}()

	queue, err := amqpCh.QueueDeclare(
		conf.AMQP.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logg.Fatal("failed to declare queue")
	}

	msgs, err := amqpCh.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
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
