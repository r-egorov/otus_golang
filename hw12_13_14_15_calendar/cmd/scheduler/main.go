package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	conn, err := amqp.Dial(conf.AMQP.Uri)
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

	store := sqlstorage.New(
		conf.Storage.User,
		conf.Storage.Password,
		conf.Storage.DBName,
		conf.Storage.Host,
		conf.Storage.Port,
	)
	err = store.Connect(ctx)
	if err != nil {
		logg.Fatal("failed to connect to DB")
	}

	//ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//body := "hello world"
	//
	//err = amqpCh.PublishWithContext(
	//	ctxTimeout,
	//	"",
	//	queue.Name,
	//	false,
	//	false,
	//	amqp.Publishing{
	//		ContentType: "text/plain",
	//		Body:        []byte(body),
	//	},
	//)
	//if err != nil {
	//	logg.Error("not published")
	//}

	logg.Info("started scheduler...")
	notifyBefore := conf.Scheduler.NotifyBefore
	go func() {
		for {
			start := time.Now().UTC()

			notifications, err := store.ListToNotify(
				ctx,
				start,
				notifyBefore,
				5*time.Minute,
			)
			fmt.Println(notifications)
			if err != nil {
				logg.Error(fmt.Sprintf("cant retrive notifications from DB: %v", err))
			}

			for _, notification := range notifications {
				jsoned, err := json.Marshal(notification)
				if err != nil {
					logg.Error(fmt.Sprintf("failed to marshal: %s", notification))
				}

				ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = amqpCh.PublishWithContext(
					ctxTimeout,
					"",
					queue.Name,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        jsoned,
					},
				)
				if err != nil {
					logg.Error(fmt.Sprintf("failed to publish: %s", notification))
				} else {
					logg.Info(fmt.Sprintf("published: %s", notification))
				}
			}

			timer := time.NewTimer(10*time.Second - time.Since(start))
			select {
			case <-timer.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()

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
