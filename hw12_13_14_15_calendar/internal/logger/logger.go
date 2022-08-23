package logger

import (
	"fmt"
	"io"
	"os"
)

type logLevel int

const (
	fatalLevel = iota
	errorLevel
	warnLevel
	infoLevel
	debugLevel
)

func (l logLevel) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}[l]
}

type Logger struct {
	level logLevel
	out   io.Writer
}

func New(out io.Writer, level string) *Logger {
	var levelCode logLevel
	switch level {
	case "FATAL":
		levelCode = fatalLevel
	case "ERROR":
		levelCode = errorLevel
	case "WARN":
		levelCode = warnLevel
	case "INFO":
		levelCode = infoLevel
	case "DEBUG":
		levelCode = debugLevel
	}
	return &Logger{
		level: levelCode,
		out:   out,
	}
}

func (l Logger) Info(msg string) {
	if l.level >= infoLevel {
		l.log(infoLevel, msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level >= errorLevel {
		l.log(errorLevel, msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level >= warnLevel {
		l.log(warnLevel, msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.level >= debugLevel {
		l.log(debugLevel, msg)
	}
}

func (l Logger) Fatal(msg string) {
	l.log(fatalLevel, msg)
	os.Exit(1)
}

func (l Logger) log(level logLevel, msg string) {
	_, err := fmt.Fprintf(l.out, "[%s] %s\n", level.String(), msg)
	if err != nil {
		panic(fmt.Errorf("fatal: logger can't write to its file, err %w", err))
	}
}
