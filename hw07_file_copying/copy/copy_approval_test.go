package copy_test

import (
	"os"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"github.com/r-egorov/otus_golang/hw07_file_copying/copy"
	"github.com/stretchr/testify/require"
)

const (
	inputFolder   = "../testdata"
	inputFilePath = "../testdata/input.txt"
)

func TestCopyApproval(t *testing.T) {
	approvalCases := []struct {
		name          string
		offset, limit int64
	}{
		{
			name:   "out offset0 limit0",
			offset: 0,
			limit:  0,
		},
		{
			name:   "out offset0 limit10",
			offset: 0,
			limit:  10,
		},
		{
			name:   "out offset0 limit1000",
			offset: 0,
			limit:  1000,
		},
		{
			name:   "out offset0 limit10000",
			offset: 0,
			limit:  10000,
		},
		{
			name:   "out offset100 limit1000",
			offset: 100,
			limit:  1000,
		},
		{
			name:   "out offset6000 limit1000",
			offset: 6000,
			limit:  1000,
		},
	}
	for _, ac := range approvalCases {
		t.Run(ac.name, func(t *testing.T) {
			inputText := getFileContent(t, inputFilePath)
			expected := inputText

			tc := setUpTestCase(inputText, expected)
			defer tc.tearDown()

			err := copy.Copy(tc.sourceFile, tc.destFile, ac.offset, ac.limit)

			require.NoError(t, err)

			gotText := getFileContent(t, tc.destFile.Name())

			approvals.UseFolder(inputFolder)
			approvals.VerifyString(t, gotText)
		})
	}
}

func getFileContent(t *testing.T, filePath string) string {
	t.Helper()
	f, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal("can't read input")
	}
	return string(f)
}
