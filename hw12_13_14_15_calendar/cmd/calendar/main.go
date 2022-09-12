package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	conf, err := config.NewConfig(configFilePath)
	if err != nil {
		log.Fatalf("config: %s", err.Error()) //nolint:gocritic
	}
	logOut, logOutClose := getLogWriter(conf)
	defer func() {
		if err := logOutClose(); err != nil {
			panic(err)
		}
	}()

	logg := logger.New(logOut, conf.Logger.Level)

	var storage app.Storage
	switch conf.Storage.StorageType {
	case config.PSQLStorageType:
		storage = sqlstorage.New(
			conf.Storage.User,
			conf.Storage.Password,
			conf.Storage.DBName,
			conf.Storage.Host,
			conf.Storage.Port,
		)
	default:
		storage = memorystorage.New()
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, conf.Server.Host, conf.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	serverStopped := make(chan struct{})
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		serverStopped <- struct{}{}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
	<-serverStopped
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
