package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("it reads the files in directory", func(t *testing.T) {
		initialEnv := map[string]string{
			"HELLO":     "hello",
			"UNSET":     "",
			"EMPTYLINE": "\nsecondline",
			"SPACES":    "spaces \t\f ",
			"ZEROTERM":  "zeroterminated\x00",
		}
		te := NewTestEnv(t)
		defer te.tearDown(t)

		te.addEnvVarFiles(t, initialEnv)

		expected := expectedEnv(initialEnv)
		got, err := ReadDir(te.tmpDirPath)

		require.NoError(t, err)
		require.Equal(t, expected, got)
	})

	t.Run("filenames can't include `=`", func(t *testing.T) {
		initialEnv := map[string]string{
			"HELLO":    "hello",
			"INVALID=": "invalid",
		}
		te := NewTestEnv(t)
		defer te.tearDown(t)

		te.addEnvVarFiles(t, initialEnv)

		_, err := ReadDir(te.tmpDirPath)

		require.Error(t, err, ErrAssignationSignInFilename)
	})
}

type TestEnv struct {
	tmpDirPath string
}

func NewTestEnv(t *testing.T) TestEnv {
	t.Helper()

	tmpDirPath, err := ioutil.TempDir("", "env_directory")
	if err != nil {
		t.Fatalf("can't create tempdir: %v", err)
	}

	return TestEnv{
		tmpDirPath: tmpDirPath,
	}
}

func (te *TestEnv) tearDown(t *testing.T) {
	t.Helper()

	err := os.RemoveAll(te.tmpDirPath)
	if err != nil {
		t.Fatalf("can't remove tempdir: %v", err)
	}
}

func (te *TestEnv) addEnvVarFile(t *testing.T, key, value string) {
	t.Helper()

	filePath := path.Join(te.tmpDirPath, key)
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		t.Fatalf("can't create env var file: %v", err)
	}
	_, err = file.Write([]byte(value))
	if err != nil {
		t.Fatalf("can't write to env var file: %v", err)
	}
}

func (te *TestEnv) addEnvVarFiles(t *testing.T, envVars map[string]string) {
	t.Helper()

	for key, value := range envVars {
		te.addEnvVarFile(t, key, value)
	}
}

func expectedEnv(env map[string]string) Environment {
	res := make(Environment)
	for key, value := range env {
		if len(value) == 0 {
			res[key] = EnvValue{
				Value:      value,
				NeedRemove: true,
			}
		} else {
			value = strings.Replace(value, "\x00", "\n", -1)
			splitted := strings.Split(value, "\n")
			firstLine := splitted[0]
			firstLine = strings.TrimRight(firstLine, " \n\t\v\f\r")
			res[key] = EnvValue{
				Value:      firstLine,
				NeedRemove: false,
			}
		}
	}
	return res
}
