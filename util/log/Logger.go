package log

import (
	"fmt"
	syslog "log"
	"os"
	"sync"
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

var lock sync.Mutex
var level = InfoLevel
var log *Logger = nil
var logPath = "./mimc.log"

func SetLogPath(path string) {
	logPath = path
}
func SetLogLevel(lvl LogLevel) {
	level = lvl
	fmt.Printf("setLogLevel: %d, %d\n", level, lvl)
}

func GetLogger() *Logger {
	if log == nil {
		lock.Lock()
		if log == nil {
			log = new(Logger)
			logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
			if nil != err {
				panic(err)
			}
			log.level = level
			fmt.Printf("GetLogger: %d, %d\n", level, level)
			log.log = syslog.New(logFile, "\r\n", syslog.Ldate|syslog.Ltime|syslog.Lshortfile)
		}
		lock.Unlock()
	}
	return log
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
