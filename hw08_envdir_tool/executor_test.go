package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	env = Environment{
		"BAR":   EnvValue{Value: "bar", NeedRemove: false},
		"EMPTY": EnvValue{Value: "", NeedRemove: true},
		"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
		"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
		"UNSET": EnvValue{Value: "", NeedRemove: true},
	}

	argsValid = []string{"testdata/echo.sh", "arg1", "arg2", "arg3"}

	scriptResult = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is ()
EMPTY is ()
arguments are arg1 arg2 arg3
`
)

func TestRunCmd(t *testing.T) {
	t.Run("Successful execution with ret-code 0", func(t *testing.T) {
		stdOut := os.Stdout

		reader, writer, err := os.Pipe()
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout = writer

		retCode := RunCmd(argsValid, env)
		require.Equal(t, 0, retCode)

		writer.Close()
		os.Stdout = stdOut

		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, reader)
		if err != nil {
			log.Fatal(err)
		}

		require.Equal(t, scriptResult, buffer.String())
	})
}
