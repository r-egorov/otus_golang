package logger

import (
	"fmt"
	"io"
)

type logLevel int

const (
	errorLevel = iota
	warnLevel
	infoLevel
	debugLevel
)

const errMsgLogfileInvalid = "logger: can't write to its outfile"

func (l logLevel) String() string {
	return [...]string{"ERROR", "WARN", "INFO", "DEBUG"}[l]
}

type Logger struct {
	level logLevel
	out   io.Writer
}

func New(out io.Writer, level string) *Logger {
	var levelCode logLevel
	switch level {
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
		l.log(msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level >= errorLevel {
		l.log(msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level >= warnLevel {
		l.log(msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.level >= debugLevel {
		l.log(msg)
	}
}

func (l Logger) log(msg string) {
	_, err := fmt.Fprintf(l.out, "%s\n", msg)
	if err != nil {
		panic(errMsgLogfileInvalid)
	}
}
