package log

import (
	"fmt"
	syslog "log"
	"os"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

type Logger struct {
	level LogLevel
	log   *syslog.Logger
}

func GetDefaultLogger() *Logger {
	return GetLogger(InfoLevel)
}
func GetLogger(level LogLevel) *Logger {
	logger := new(Logger)
	logger.log = syslog.New(os.Stdout, "\r\n", syslog.Ldate|syslog.Ltime|syslog.Lshortfile)
	logger.level = level
	return logger
}

func (this *Logger) Info(format string, args ...interface{}) {
	if this.level <= InfoLevel {
		this.log.SetPrefix("[info] ")
		this.log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (this *Logger) Debug(format string, args ...interface{}) {
	if this.level <= DebugLevel {
		this.log.SetPrefix("[debug]")
		this.log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (this *Logger) Warn(format string, args ...interface{}) {
	if this.level <= WarnLevel {
		this.log.SetPrefix("[warn]")
		this.log.Output(2, fmt.Sprintf(format, args...))
	}
}

func (this *Logger) Error(format string, args ...interface{}) {
	if this.level <= ErrorLevel {
		this.log.SetPrefix("[error]")
		this.log.Output(2, fmt.Sprintf(format, args...))
	}
}
