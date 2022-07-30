package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
	dirFile, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer dirFile.Close()

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, f := range dirEntries {
		envFileName := f.Name()
		if skipEnvFile(envFileName) {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", dir, envFileName)
		value, err := extractValueFromFile(filePath)
		if err != nil {
			continue
		}

		envValue := EnvValue{Value: value, NeedRemove: len(value) == 0}
		env[envFileName] = envValue
	}

	return env, nil
}

func skipEnvFile(fileName string) bool {
	return strings.Contains(fileName, "=") || strings.HasPrefix(fileName, ".")
}

func extractValueFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", nil
	}
	value = purgeEnvValue(value)
	return value, nil
}

func purgeEnvValue(value string) string {
	value = strings.TrimSuffix(value, "\n")
	value = strings.TrimSuffix(value, " ")
	value = strings.ReplaceAll(value, "\u0000", "\n")
	return value
}
