package copy

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const defaultBufferSize = 512

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type files struct {
	sourceFd, destFd *os.File
}

func Copy(sourceFd, destFd *os.File, offset, limit int64) error {
	lenToCopy, err := calculateLenToCopy(sourceFd, offset, limit)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(lenToCopy)
	barReader := bar.NewProxyReader(sourceFd)

	err = copyContent(barReader, destFd, lenToCopy)
	if err != nil {
		return err
	}

	bar.Finish()
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
		lenToCopy = limit
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
