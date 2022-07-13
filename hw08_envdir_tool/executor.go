package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

const failCode = -1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(
	cmd []string, env Environment,
	stdout, stderr io.Writer, stdin io.Reader,
) (returnCode int) {
	var command *exec.Cmd

	if len(cmd) < 1 {
		return 125 // FIXME
	} else if len(cmd) == 1 {
		command = exec.Command(cmd[0])
	} else {
		command = exec.Command(cmd[0], cmd[1:]...)
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
		if _, ok := err.(*exec.ExitError); !ok {
			return failCode
		}
	}

	return command.ProcessState.ExitCode()
}

func prepareEnvironment(env Environment) ([]string, error) {
	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return nil, err
			}
		} else {
			err := os.Setenv(key, value.Value)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Environ(), nil
}
