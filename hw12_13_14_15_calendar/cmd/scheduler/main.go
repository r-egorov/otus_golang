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
	"time"

	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/rmq"
	sqlstorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

const scanPeriod = 5 * time.Minute

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

	logg.Info("started scheduler...")
	notifyBefore := conf.Scheduler.NotifyBefore
	go func() {
		for {
			start := time.Now().UTC()

			notifications, err := store.ListToNotify(
				ctx,
				start,
				notifyBefore,
				scanPeriod,
			)
			fmt.Println(notifications)
			if err != nil {
				logg.Error(fmt.Sprintf("cant retrieve notifications from DB: %v", err))
			}

			for _, notification := range notifications {
				jsoned, err := json.Marshal(notification)
				if err != nil {
					logg.Error(fmt.Sprintf("failed to marshal: %s", notification))
				}

				ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err = rabbit.Publish(ctxTimeout, jsoned)
				if err != nil {
					logg.Error(fmt.Sprintf("failed to publish: %s", notification))
				} else {
					logg.Info(fmt.Sprintf("published: %s", notification))
				}
			}

			err = store.ClearOlderThanYear(ctx)
			if err != nil {
				logg.Error(fmt.Sprintf("failed to delete old events: %v", err))
			}

			timer := time.NewTimer(conf.Scheduler.Period - time.Since(start))
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
			log.Fatalf("fatal: log file %s, err: %w", c.Logger.OutPath, err)
		}
	}
	outClose = func() error { return out.Close() }
	return
}
