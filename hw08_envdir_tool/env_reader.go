package main

import (
	"errors"
	"os"
	"path"
	"strings"
)

const (
	whiteSpaces    = " \n\t\v\f\r"
	zeroTerminator = "\x00"
)

var (
	ErrAssignationSignInFilename = errors.New("filename can't contain assignation sign")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make(Environment, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		envValue, err := createEnvValueFromFile(dir, file)
		if err != nil {
			return nil, err
		}

		res[file.Name()] = envValue
	}
	return res, nil
}

func createEnvValueFromFile(dir string, file os.DirEntry) (EnvValue, error) {
	fileInfo, err := file.Info()
	if err != nil {
		return EnvValue{}, err
	}

	if strings.Contains(fileInfo.Name(), "=") {
		return EnvValue{}, ErrAssignationSignInFilename
	}

	preparedValue := ""
	needRemove := false

	if fileInfo.Size() == 0 {
		needRemove = true
	} else {
		content, err := os.ReadFile(path.Join(dir, fileInfo.Name()))
		if err != nil {
			return EnvValue{}, err
		}
		preparedValue = prepareValue(string(content))
	}

	return EnvValue{
		Value:      preparedValue,
		NeedRemove: needRemove,
	}, nil
}

func prepareValue(content string) string {
	content = strings.Replace(content, zeroTerminator, "\n", -1)
	splitted := strings.Split(content, "\n")
	firstLine := splitted[0]
	firstLine = strings.TrimRight(firstLine, whiteSpaces)

	return firstLine
}
