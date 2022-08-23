package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server/http"
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

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.NewConfig(configFilePath)
	if err != nil {
		log.Fatalf("config: %s", err.Error()) //nolint:gocritic
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

	// Configure Application
	calendar := app.New(logg, conf)
	if err = calendar.ConnectToStorage(ctx); err != nil {
		logg.Fatal("can't connect to the storage")
	}
	defer func() {
		err := calendar.DisconnectFromStorage(ctx)
		if err != nil {
			logg.Fatal("can't disconnect from the storage")
		}
	}()

	// Configure API servers
	httpserver := internalhttp.NewServer(logg, calendar, conf.Server.Host, conf.Server.Port)
	grpcserver := internalgrpc.NewService(logg, calendar, conf.Server.Host, "9000") // FIXME

	httpserver.Start(ctx)
	if err := grpcserver.Start(ctx); err != nil {
		logg.Fatal("failed to start grpc server: " + err.Error())
	}

	logg.Info("calendar is running...")

	<-ctx.Done()

	ctxStop, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := httpserver.Stop(ctxStop); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
	grpcserver.Stop(ctxStop)
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
