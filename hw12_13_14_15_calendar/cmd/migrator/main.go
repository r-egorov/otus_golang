package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/config"
)

var (
	flags          = flag.NewFlagSet("goose", flag.ExitOnError)
	dir            = flags.String("dir", ".", "directory with migration files")
	configFilePath = flags.String("config", "/etc/calendar/config.toml", "Path to configuration file")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	conf, err := config.NewConfig(*configFilePath)
	if err != nil {
		log.Fatalf("config: %s", err.Error())
	}

	if conf.Storage.StorageType == config.InmemoryStorageType {
		log.Fatalf("migrator: inmemory storage in use, can't migrate")
	}
	dbstring := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Storage.User,
		conf.Storage.Password,
		conf.Storage.Host,
		conf.Storage.Port,
		conf.Storage.DBName,
	)

	db, err := goose.OpenDBWithDriver("postgres", dbstring)
	if err != nil {
		log.Fatalf("migrator: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("migrator: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(args[0], db, *dir, arguments...); err != nil {
		log.Fatalf("migrator %v: %v", args[0], err) //nolint:gocritic
	}
}
