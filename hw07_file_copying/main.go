package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/r-egorov/otus_golang/hw07_file_copying/copy"
)

var (
	from, to              string
	limit, offset         int64
	ErrSourceFileNotFound = errors.New("source file not found")
	ErrDestFileInvalid    = errors.New("can't create or open dest file")
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	// Prepare source FD
	sourceFd, err := os.OpenFile(from, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalf("%s: %s", ErrSourceFileNotFound, from)
		return
	}
	defer sourceFd.Close()

	// Prepare dest FD
	destFd, err := os.Create(to)
	if err != nil {
		log.Fatal(ErrDestFileInvalid)
		return
	}
	defer destFd.Close()

	copy.Copy(sourceFd, destFd, offset, limit)
}
