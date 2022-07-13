package main

import (
	"log"
	"os"
)

const usageDirections = "Usage: <go-envdir> <envdirectory> <cmd> [arg...]"

func main() {
	if len(os.Args) < 3 {
		log.Fatal(usageDirections)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	os.Exit(RunCmd(os.Args[2:], env, os.Stdout, os.Stderr, os.Stdin))
}
