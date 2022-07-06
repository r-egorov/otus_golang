package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	tmpDir        = "tmp"
	tmpSourcename = "tmp.txt"
	tmpDestname   = "tmpdest.txt"
	testText      = "sometext"
)

type testCase struct {
	sourceFile, destFile     *os.File
	sourceText, expectedText string
}

func (t *testCase) tearDown() {
	err := os.Remove(t.sourceFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(t.destFile.Name())
	if err != nil {
		log.Fatal(err)
	}
}

func setUpTestCase(sourceText, expectedText string) testCase {
	tmpSourceFile, err := os.CreateTemp("", tmpSourcename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tmpSourceFile.Write([]byte(sourceText))
	if err != nil {
		log.Fatal(err)
	}

	tmpDestFile, err := os.CreateTemp("", tmpDestname)
	if err != nil {
		log.Fatal(err)
	}

	return testCase{
		sourceFile:   tmpSourceFile,
		destFile:     tmpDestFile,
		sourceText:   sourceText,
		expectedText: expectedText,
	}
}

func TestCopyContent(t *testing.T) {
	t.Run("copies content of a reader to writer", func(t *testing.T) {
		expected := testText
		lenToCopy := int64(len(testText))

		source := &bytes.Buffer{}
		source.WriteString(testText)

		dest := &bytes.Buffer{}
		err := copyContent(source, dest, lenToCopy)

		require.NoError(t, err)

		got := dest.String()

		require.Equal(t, expected, got)
	})
}

func TestCopySuccess(t *testing.T) {
	t.Run("copies from one file to another", func(t *testing.T) {
		tc := setUpTestCase(testText, testText)

		defer tc.tearDown()

		err := Copy(tc.sourceFile.Name(), tc.destFile.Name(), 0, 0)

		require.NoError(t, err)

		body, err := os.ReadFile(tc.destFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		got := string(body)
		require.Equal(t, tc.expectedText, got)
	})

	t.Run("copies with offset", func(t *testing.T) {
		offset = 5
		expected := testText[offset:]

		tc := setUpTestCase(testText, expected)
		defer tc.tearDown()

		err := Copy(tc.sourceFile.Name(), tc.destFile.Name(), offset, 0)

		require.NoError(t, err)

		body, err := os.ReadFile(tc.destFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		got := string(body)
		require.Equal(t, tc.expectedText, got)
	})
}

func TestCopyFail(t *testing.T) {
	t.Run("no source file", func(t *testing.T) {
		tc := setUpTestCase(testText, "")
		defer tc.tearDown()

		err := Copy("invalidpath", tc.destFile.Name(), 0, 0)
		require.Error(t, err, ErrSourceFileNotFound)
	})

	t.Run("offset is greater than the source file length", func(t *testing.T) {
		var offset int64 = 9999999

		tc := setUpTestCase(testText, "")
		defer tc.tearDown()

		err := Copy(tc.sourceFile.Name(), tc.destFile.Name(), offset, 0)

		require.Error(t, err, ErrOffsetExceedsFileSize)
	})
}
