package mycopy_test

import (
	"os"
	"testing"

	"github.com/r-egorov/otus_golang/hw07_file_copying/mycopy"
	"github.com/stretchr/testify/require"
)

const (
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

func TestCopySuccess(t *testing.T) {
	t.Run("copies from one file to another", func(t *testing.T) {
		te := setUpTestEnv(t, testText)
		defer te.tearDown(t)

		expected := testText
		err := mycopy.Copy(te.sourceFile, te.destFile, 0, 0)

		require.NoError(t, err)

		got := getFileContent(t, te.destFile.Name())
		require.Equal(t, expected, got)
	})

	t.Run("copies with offset", func(t *testing.T) {
		offset := int64(50)
		expected := testText[offset:]

		te := setUpTestEnv(t, testText)
		defer te.tearDown(t)

		err := mycopy.Copy(te.sourceFile, te.destFile, offset, 0)

		require.NoError(t, err)

		got := getFileContent(t, te.destFile.Name())
		require.Equal(t, expected, got)
	})

	t.Run("copies with limit", func(t *testing.T) {
		limit := int64(50)
		expected := testText[:limit]

		te := setUpTestEnv(t, testText)
		defer te.tearDown(t)

		err := mycopy.Copy(te.sourceFile, te.destFile, 0, limit)

		require.NoError(t, err)

		got := getFileContent(t, te.destFile.Name())
		require.Equal(t, expected, got)
	})

	t.Run("copies with offset and limit", func(t *testing.T) {
		offset := int64(50)
		limit := int64(50)
		expected := testText[offset : offset+limit]

		te := setUpTestEnv(t, testText)
		defer te.tearDown(t)

		err := mycopy.Copy(te.sourceFile, te.destFile, offset, limit)

		require.NoError(t, err)

		got := getFileContent(t, te.destFile.Name())
		require.Equal(t, expected, got)
	})
}

func TestCopyFail(t *testing.T) {
	t.Run("offset is greater than the source file length", func(t *testing.T) {
		var offset int64 = 9999999

		tc := setUpTestEnv(t, testText)
		defer tc.tearDown(t)

		err := mycopy.Copy(tc.sourceFile, tc.destFile, offset, 0)

		require.Error(t, err, mycopy.ErrOffsetExceedsFileSize)
	})

	t.Run("does not copy from /dev/urandom", func(t *testing.T) {
		te := setUpTestEnv(t, testText)
		defer te.tearDown(t)

		rand, err := os.Open("/dev/urandom")
		if err != nil {
			t.Fatal("can't open /dev/urandom")
		}

		err = mycopy.Copy(rand, te.destFile, 0, 0)
		require.Error(t, err, mycopy.ErrUnsupportedFile)
	})
}
