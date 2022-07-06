package copy_test

import (
	"log"
	"os"
	"testing"

	"github.com/r-egorov/otus_golang/hw07_file_copying/copy"
	"github.com/stretchr/testify/require"
)

const (
	tmpDir        = "tmp"
	tmpSourcename = "tmp.txt"
	tmpDestname   = "tmpdest.txt"
	testText      = `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
					sed do eiusmod tempor incididunt ut labore et dolore magna 
					aliqua. Tortor dignissim convallis aenean et tortor at. 
					Lorem dolor sed viverra ipsum nunc aliquet. Est lorem ipsum 
					dolor sit amet consectetur adipiscing elit. Vestibulum lectus 
					mauris ultrices eros in cursus turpis massa tincidunt. 
					Pellentesque adipiscing commodo elit at imperdiet dui. 
					Cursus turpis massa tincidunt dui ut ornare. Massa vitae 
					tortor condimentum lacinia quis vel. Id diam maecenas ultricies 
					mi eget mauris pharetra et ultrices. Porttitor eget dolor morbi 
					non. Arcu risus quis varius quam quisque id diam vel. Eget nunc 
					scelerisque viverra mauris in aliquam sem fringilla ut. Leo a 
					diam sollicitudin tempor id eu nisl. Feugiat vivamus at augue 
					eget arcu dictum varius duis. Natoque penatibus et magnis dis. 
					Massa id neque aliquam vestibulum. Morbi non arcu risus quis 
					varius quam quisque id diam. Eu turpis egestas pretium aenean 
					pharetra. Posuere sollicitudin aliquam ultrices sagittis orci a.`
)

type testCase struct {
	sourceFile, destFile     *os.File
	sourceText, expectedText string
}

func (t *testCase) tearDown() {
	os.Remove(t.sourceFile.Name())
	os.Remove(t.destFile.Name())
}

func setUpTestCase(sourceText, expectedText string) testCase {
	tmpSourceFile, err := os.CreateTemp("", tmpSourcename)
	if err != nil {
		log.Fatal(err)
	}
	tmpSourceFile.Write([]byte(sourceText))

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

func TestCopySuccess(t *testing.T) {
	t.Run("copies from one file to another", func(t *testing.T) {
		tc := setUpTestCase(testText, testText)

		defer tc.tearDown()

		err := copy.Copy(tc.sourceFile.Name(), tc.destFile.Name(), 0, 0)

		require.NoError(t, err)

		body, err := os.ReadFile(tc.destFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		got := string(body)
		require.Equal(t, tc.expectedText, got)
	})

	t.Run("copies with offset", func(t *testing.T) {
		offset := int64(50)
		expected := testText[offset:]

		tc := setUpTestCase(testText, expected)
		defer tc.tearDown()

		err := copy.Copy(tc.sourceFile.Name(), tc.destFile.Name(), offset, 0)

		require.NoError(t, err)

		body, err := os.ReadFile(tc.destFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		got := string(body)
		require.Equal(t, tc.expectedText, got)
	})

	t.Run("copies with limit", func(t *testing.T) {
		limit := int64(50)
		expected := testText[:limit+1]

		tc := setUpTestCase(testText, expected)
		defer tc.tearDown()

		err := copy.Copy(tc.sourceFile.Name(), tc.destFile.Name(), 0, limit)

		require.NoError(t, err)

		body, err := os.ReadFile(tc.destFile.Name())
		if err != nil {
			log.Fatal(err)
		}

		got := string(body)
		require.Equal(t, tc.expectedText, got)
	})

	t.Run("copies with offset and limit", func(t *testing.T) {
		offset := int64(50)
		limit := int64(50)
		expected := testText[offset : offset+limit+1]

		tc := setUpTestCase(testText, expected)
		defer tc.tearDown()

		err := copy.Copy(tc.sourceFile.Name(), tc.destFile.Name(), offset, limit)

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

		err := copy.Copy("invalidpath", tc.destFile.Name(), 0, 0)
		require.Error(t, err, copy.ErrSourceFileNotFound)
	})

	t.Run("offset is greater than the source file length", func(t *testing.T) {
		var offset int64 = 9999999

		tc := setUpTestCase(testText, "")
		defer tc.tearDown()

		err := copy.Copy(tc.sourceFile.Name(), tc.destFile.Name(), offset, 0)

		require.Error(t, err, copy.ErrOffsetExceedsFileSize)
	})
}
