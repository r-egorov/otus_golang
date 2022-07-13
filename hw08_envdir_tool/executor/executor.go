package executor

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/r-egorov/otus_golang/hw08_envdir_tool/envreader"
)

const failCode = -1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(
	cmd []string, env envreader.Environment,
	stdout, stderr io.Writer, stdin io.Reader,
) (returnCode int) {
	var command *exec.Cmd

	switch {
	case len(cmd) == 1:
		command = exec.Command(cmd[0]) //nolint:gosec
	case len(cmd) > 1:
		command = exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	default:
		return failCode
	}

	command.Stderr = stderr
	command.Stdin = stdin
	command.Stdout = stdout

	preparedEnv, err := prepareEnvironment(env)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return failCode
	}

	command.Env = preparedEnv
	err = command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			return failCode
		}
	}

	return command.ProcessState.ExitCode()
}

func prepareEnvironment(env envreader.Environment) ([]string, error) {
	for key, value := range env {
		var err error
		if value.NeedRemove {
			err = os.Unsetenv(key)
		} else {
			err = os.Setenv(key, value.Value)
		}
		if err != nil {
			return nil, err
		}
	}
	return os.Environ(), nil
}
