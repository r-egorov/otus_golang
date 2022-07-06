package mycopy_test

import (
	"os"
	"testing"
)

type testEnv struct {
	sourceFile, destFile *os.File
	sourceText           string
}

func (t *testEnv) tearDown() {
	t.sourceFile.Close()
	t.destFile.Close()
	os.Remove(t.sourceFile.Name())
	os.Remove(t.destFile.Name())
}

func setUpTestEnv(t *testing.T, sourceText string) testEnv {
	t.Helper()
	tmpSourceFile, err := os.CreateTemp("", tmpSourcename)
	if err != nil {
		t.Fatal(err)
	}
	tmpSourceFile.Write([]byte(sourceText))

	tmpDestFile, err := os.CreateTemp("", tmpDestname)
	if err != nil {
		t.Fatal(err)
	}

	return testEnv{
		sourceFile: tmpSourceFile,
		destFile:   tmpDestFile,
		sourceText: sourceText,
	}
}

func getFileContent(t *testing.T, filePath string) string {
	t.Helper()
	f, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("can't read input %s", filePath)
	}
	return string(f)
}
