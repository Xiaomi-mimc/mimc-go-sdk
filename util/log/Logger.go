package log

import (
	"fmt"
	syslog "log"
	"os"
	"strconv"
	"sync"
	"time"
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
		fmt.Printf("["+currTime()+"] "+format+"\n", args...)
	}
}

func (this *Logger) Debug(format string, args ...interface{}) {
	if this.level <= DebugLevel {
		this.log.SetPrefix("[debug]")
		this.log.Output(2, fmt.Sprintf(format, args...))
		fmt.Printf("["+currTime()+"] "+format+"\n", args...)
	}
}

func (this *Logger) Warn(format string, args ...interface{}) {
	if this.level <= WarnLevel {
		this.log.SetPrefix("[warn]")
		this.log.Output(2, fmt.Sprintf(format, args...))
		fmt.Printf("["+currTime()+"] "+format+"\n", args...)
	}
}

func (this *Logger) Error(format string, args ...interface{}) {
	if this.level <= ErrorLevel {
		this.log.SetPrefix("[error]")
		this.log.Output(2, fmt.Sprintf(format, args...))
		fmt.Printf("["+currTime()+"] "+format+"\n", args...)
	}
}

func currTime() string {
	now := time.Now()
	const base_format = "2006-01-02 15:04:05"
	nsecStr := strconv.Itoa(now.Nanosecond())
	timeStr := now.Format(base_format) + "." + nsecStr[0:4]
	return timeStr
}
