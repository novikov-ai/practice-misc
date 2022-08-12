package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	TestDataPath    = "testdata/"
	TestDataEnvPath = "testdata/env"
)

func TestReadDir(t *testing.T) {
	t.Run("Correct environment was configured", func(t *testing.T) {
		envExpected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		envActual, err := ReadDir(TestDataEnvPath)
		require.Nil(t, err)

		require.Equal(t, envExpected, envActual)
	})

	t.Run("Directory is empty or contains invalid files", func(t *testing.T) {
		tmpDirPath, err := os.MkdirTemp(TestDataPath, "env_test")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDirPath)

		tmpDirFile, err := os.CreateTemp(tmpDirPath, "=")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpDirFile.Name())

		env, err := ReadDir(tmpDirPath)
		require.Nil(t, err)

		require.Equal(t, 0, len(env))
	})
}
