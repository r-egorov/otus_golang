package mycopy

import (
	"errors"
	"io"
	"os"

	"github.com/r-egorov/otus_golang/hw07_file_copying/progressbar"
)

const defaultBufferSize = 4096

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy takes source and destination file descriptors,
// copies at most `limit` bytes from source to dest,
// starting from the `offset` byte.
// In case `limit` is zero, copies the content till the EOF
// Outputs the progress bar while copying.
func Copy(sourceFd, destFd *os.File, offset, limit int64) error {
	lenToCopy, err := calculateLenToCopy(sourceFd, offset, limit)
	if err != nil {
		return err
	}

	bar := progressbar.NewBar(lenToCopy)
	barReader := bar.NewProxyReader(sourceFd)

	bar.Start()
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

	if !fileStat.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
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
