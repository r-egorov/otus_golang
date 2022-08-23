package main

import (
	"flag"
	"fmt"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/logger"
	"log"
	"os"
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

	logOut, logOutClose := getLogWriter(conf)
	defer func() {
		if err := logOutClose(); err != nil {
			panic(err)
		}
	}()
	logg := logger.New(logOut, conf.Logger.Level)

	calendar := app.New(logg, conf)
	calendar.Run()
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
