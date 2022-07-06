package mycopy_test

import (
	"os"
	"testing"
)

type testEnv struct {
	sourceFile, destFile *os.File
	sourceText           string
}

func (te *testEnv) tearDown(t *testing.T) {
	t.Helper()
	err := te.sourceFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = te.destFile.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(te.sourceFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(te.destFile.Name())
	if err != nil {
		t.Fatal(err)
	}
}

func setUpTestEnv(t *testing.T, sourceText string) testEnv {
	t.Helper()
	tmpSourceFile, err := os.CreateTemp("", tmpSourcename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpSourceFile.Write([]byte(sourceText))
	if err != nil {
		t.Fatal(err)
	}

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
