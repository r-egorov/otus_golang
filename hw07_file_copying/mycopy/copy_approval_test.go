package mycopy_test

import (
	"path"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/r-egorov/otus_golang/hw07_file_copying/mycopy"
	"github.com/stretchr/testify/require"
)

var (
	inputFilePath      = path.Join("..", "testdata", "input.txt")
	emptyInputFilePath = path.Join("..", "testdata", "empty_input.txt")
	inputFolder        = path.Dir(inputFilePath)
)

func TestCopyApproval(t *testing.T) {
	testCases := []struct {
		name, filePath string
		offset, limit  int64
	}{
		{
			name:     "out offset0 limit0",
			filePath: inputFilePath,
			offset:   0,
			limit:    0,
		},
		{
			name:     "out offset0 limit10",
			filePath: inputFilePath,
			offset:   0,
			limit:    10,
		},
		{
			name:     "out offset0 limit1000",
			filePath: inputFilePath,
			offset:   0,
			limit:    1000,
		},
		{
			name:     "out offset0 limit10000",
			filePath: inputFilePath,
			offset:   0,
			limit:    10000,
		},
		{
			name:     "out offset100 limit1000",
			filePath: inputFilePath,
			offset:   100,
			limit:    1000,
		},
		{
			name:     "out offset6000 limit1000",
			filePath: inputFilePath,
			offset:   6000,
			limit:    1000,
		},
		{
			name:     "empty offset0 limit0",
			filePath: emptyInputFilePath,
			offset:   0,
			limit:    0,
		},
		{
			name:     "empty offset0 limit10",
			filePath: emptyInputFilePath,
			offset:   0,
			limit:    10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputText := getFileContent(t, tc.filePath)

			te := setUpTestEnv(t, inputText)
			defer te.tearDown(t)

			err := mycopy.Copy(te.sourceFile, te.destFile, tc.offset, tc.limit)

			require.NoError(t, err)

			gotText := getFileContent(t, te.destFile.Name())

			approvals.UseFolder(inputFolder)
			approvals.VerifyString(t, gotText)
		})
	}
}
