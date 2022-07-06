package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/r-egorov/otus_golang/hw07_file_copying/mycopy"
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
	sourceFd, err := os.OpenFile(from, os.O_RDONLY, 0o0755)
	if err != nil {
		fmt.Printf("Error: %s: %s\n", ErrSourceFileNotFound, from)
		return
	}
	defer sourceFd.Close()

	// Prepare dest FD
	destFd, err := os.Create(to)
	if err != nil {
		fmt.Printf("Error: %s: %s", ErrDestFileInvalid, to)
		return
	}
	defer destFd.Close()

	err = mycopy.Copy(sourceFd, destFd, offset, limit)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
