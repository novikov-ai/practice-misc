package logger

import (
	"fmt"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"io"
	"log"
	"os"
	"strings"
)

type Logger struct {
	LogLevel int
	LogFile  string
	Writer   io.Writer
}

var logLevels = map[string]int{
	"DEBUG": Debug,
	"INFO":  Info,
	"WARN":  Warning,
	"ERROR": Error,
}

const (
	Debug = iota
	Info
	Warning
	Error
)

func New(conf configs.Config) *Logger {
	level := strings.TrimSpace(strings.ToUpper(conf.Logger.Level))
	lvl, ok := logLevels[level]
	if !ok {
		log.Fatal("incorrect logger level")
	}

	return &Logger{LogLevel: lvl, Writer: os.Stdout}
}

func (l Logger) log(head, message string) {
	bytes := []byte(fmt.Sprintf("[%s]: %s\n", head, message))
	_, err := l.Writer.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func (l Logger) Debug(msg string) {
	if l.LogLevel > Debug {
		return
	}

	l.log("DEBUG", msg)
}

func (l Logger) Info(msg string) {
	if l.LogLevel > Info {
		return
	}

	l.log("INFO", msg)
}

func (l Logger) Warn(msg string) {
	if l.LogLevel > Warning {
		return
	}

	l.log("WARN", msg)
}

func (l Logger) Error(msg string) {
	l.log("ERROR", msg)
}
