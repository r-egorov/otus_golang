package executor_test

import (
	"bytes"
	"github.com/r-egorov/otus_golang/hw08_envdir_tool/envreader"
	"github.com/r-egorov/otus_golang/hw08_envdir_tool/executor"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	catTestScriptPath    = path.Join("..", "testdata", "cat.sh")
	stderrTestScriptPath = path.Join("..", "testdata", "send_to_stderr.sh")
	exitTestScriptPath   = path.Join("..", "testdata", "exitcode.sh")
	echoTestScriptPath   = path.Join("..", "testdata", "echo.sh")
)

const exitCodeInTestFile = 42

func TestRunCmd(t *testing.T) {
	t.Run("it runs echo with env", func(t *testing.T) {
		out := &bytes.Buffer{}

		err := os.Setenv("ADDED", "from original env")
		if err != nil {
			t.Fatal("can't set env variable")
		}
		defer os.Unsetenv("ADDED")

		env := envreader.Environment{
			"BAR": envreader.EnvValue{
				"bar", false,
			},
			"EMPTY": envreader.EnvValue{
				"", false,
			},
			"FOO": envreader.EnvValue{
				"   foo\nwith new line", false,
			},
			"HELLO": envreader.EnvValue{
				`"hello"`, false,
			},
			"UNSET": envreader.EnvValue{
				"", true,
			},
		}
		cmd := []string{"bash", echoTestScriptPath, "arg1=1", "arg2=2"}
		exitCode := executor.RunCmd(cmd, env, out, nil, nil)

		got := out.String()
		expected := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`
		require.Equal(t, 0, exitCode)
		require.Equal(t, expected, got)
	})

	t.Run("it redirects the stdin", func(t *testing.T) {
		toSendToIn := "This goes into IN"
		in := bytes.NewBuffer([]byte(toSendToIn))
		out := &bytes.Buffer{}

		cmd := []string{"bash", catTestScriptPath}
		exitCode := executor.RunCmd(cmd, nil, out, nil, in)

		got := out.String()
		expected := toSendToIn

		require.Equal(t, 0, exitCode)
		require.Equal(t, expected, got)
	})

	t.Run("exit code equals cmd's exit code", func(t *testing.T) {
		cmd := []string{"bash", exitTestScriptPath}
		exitCode := executor.RunCmd(cmd, nil, nil, nil, nil)

		require.Equal(t, exitCodeInTestFile, exitCode)
	})

	t.Run("it redirects the stderr", func(t *testing.T) {
		toSendToIn := "This will be redirected to ERR"
		in := bytes.NewBuffer([]byte(toSendToIn))
		err := &bytes.Buffer{}

		cmd := []string{"bash", stderrTestScriptPath}
		exitCode := executor.RunCmd(cmd, nil, nil, err, in)

		got := err.String()
		expected := toSendToIn

		require.Equal(t, 0, exitCode)
		require.Equal(t, expected, got)
	})
}
