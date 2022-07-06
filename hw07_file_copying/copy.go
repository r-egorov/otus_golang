package main

import (
	"errors"
	"io"
	"os"
)

const bufferSize = 1

var (
	ErrSourceFileNotFound    = errors.New("source file not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Prepare source FD
	sourceFd, err := os.OpenFile(fromPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	fileStat, err := sourceFd.Stat()
	if err != nil {
		return err
	}

	if offset > fileStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	lenToCopy := fileStat.Size() - offset
	if limit > 0 && lenToCopy > limit {
		lenToCopy = limit
	}

	_, err = sourceFd.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Prepare dest FD
	destFd, err := os.Create(toPath)
	if err != nil {
		return err
	}

	err = copyContent(sourceFd, destFd, lenToCopy)
	if err != nil {
		return err
	}
	return nil
}

func copyContent(source io.Reader, dest io.Writer, lenToCopy int64) error {
	var haveReadBytes int64
	for haveReadBytes < lenToCopy {
		readBytes, err := io.CopyN(dest, source, bufferSize)
		if err != nil {
			return err
		}
		haveReadBytes += readBytes
	}
	return nil
}
