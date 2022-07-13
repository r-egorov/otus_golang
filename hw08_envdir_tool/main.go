package main

import (
	"github.com/r-egorov/otus_golang/hw08_envdir_tool/envreader"
	"github.com/r-egorov/otus_golang/hw08_envdir_tool/executor"
	"log"
	"os"
)

const usageDirections = "Usage: <go-envdir> <envdirectory> <cmd> [arg...]"

func main() {
	if len(os.Args) < 3 {
		log.Fatal(usageDirections)
	}
	env, err := envreader.ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	os.Exit(executor.RunCmd(os.Args[2:], env, os.Stdout, os.Stderr, os.Stdin))
}
