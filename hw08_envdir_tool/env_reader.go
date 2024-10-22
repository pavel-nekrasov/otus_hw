package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

var errWrongName = errors.New("wrong file name")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(map[string]EnvValue)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if strings.Contains(fileName, "=") {
			return nil, errWrongName
		}

		value, needRemove, err := readEnvValue(path.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		envs[fileName] = EnvValue{
			Value:      value,
			NeedRemove: needRemove,
		}
	}
	return envs, nil
}

func readEnvValue(path string) (string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", false, err
	}

	if fileInfo.Size() == 0 {
		return "", true, nil
	}

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", false, err
	}

	data = bytes.ReplaceAll(data, []byte{0}, []byte{'\n'})
	data = bytes.TrimRight(data, "\r\t \n")

	return string(data), false, nil
}
