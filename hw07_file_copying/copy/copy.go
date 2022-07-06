package copy

import (
	"errors"
	"io"
	"os"
)

const defaultBufferSize = 512

var (
	ErrSourceFileNotFound    = errors.New("source file not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type files struct {
	sourceFd, destFd *os.File
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Prepare source FD
	sourceFd, err := os.OpenFile(fromPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	lenToCopy, err := calculateLenToCopy(sourceFd, offset, limit)
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

func calculateLenToCopy(sourceFd *os.File, offset, limit int64) (int64, error) {
	fileStat, err := sourceFd.Stat()
	if err != nil {
		return 0, err
	}

	if offset > fileStat.Size() {
		return 0, ErrOffsetExceedsFileSize
	}
	offset, err = sourceFd.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	lenToCopy := fileStat.Size() - offset
	if limit > 0 && lenToCopy > limit {
		lenToCopy = limit + 1
	}
	return lenToCopy, nil
}

func copyContent(source io.Reader, dest io.Writer, lenToCopy int64) error {
	var totalReadBytes int64
	bufferSize := int64(defaultBufferSize)

	for totalReadBytes < lenToCopy {
		if bufferSize > lenToCopy-totalReadBytes {
			bufferSize = lenToCopy - totalReadBytes
		}

		readBytes, err := io.CopyN(dest, source, bufferSize)
		if err != nil {
			return err
		}
		totalReadBytes += readBytes
	}
	return nil
}
