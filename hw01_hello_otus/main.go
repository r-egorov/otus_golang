package main

import (
	"log"
	"os"

	"github.com/r-egorov/otus_golang/hw01_hello_otus/rprint"
)

func main() {
	err := rprint.RevPrint(os.Stdout, "Hello, OTUS!")
	if err != nil {
		log.Fatalf("error when output %v", err)
	}
}
