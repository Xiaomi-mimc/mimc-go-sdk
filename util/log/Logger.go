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
var level LogLevel
var log *Logger
var logPath string
var zero = [...]string{"0000", "000", "00", "0"}

func init() {
	level = InfoLevel
	log = nil
	logPath = "./mimc.log"
}

func SetLogPath(path string) {
	logPath = path
}
func SetLogLevel(lvl LogLevel) {
	if log != nil {
		log.level = lvl
	} else {
		level = lvl
	}
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
		this.log.SetPrefix("[info]\t")
		this.log.Output(2, fmt.Sprintf("\t\t"+format, args...))
		fmt.Printf("["+currTime()+"]  [info] "+format+"\n", args...)
	}
}

func (this *Logger) Debug(format string, args ...interface{}) {
	if this.level <= DebugLevel {
		this.log.SetPrefix("[debug]\t")
		this.log.Output(2, fmt.Sprintf("\t\t"+format, args...))
		fmt.Printf("["+currTime()+"] [debug] "+format+"\n", args...)
	}
}

func (this *Logger) Warn(format string, args ...interface{}) {
	if this.level <= WarnLevel {
		this.log.SetPrefix("[warn]\t")
		this.log.Output(2, fmt.Sprintf("\t\t"+format, args...))
		fmt.Printf("["+currTime()+"]  [warn] "+format+"\n", args...)
	}
}

func (this *Logger) Error(format string, args ...interface{}) {
	if this.level <= ErrorLevel {
		this.log.SetPrefix("[error]\t")
		this.log.Output(2, fmt.Sprintf("\t\t"+format, args...))
		fmt.Printf("["+currTime()+"] [error] "+format+"\n", args...)
	}
}

func currTime() string {
	now := time.Now()
	const base_format = "2006-01-02 15:04:05"
	nsecStr := strconv.Itoa(now.Nanosecond())
	nsecLen := len(nsecStr)
	if nsecLen < 4 {
		nsecStr = nsecStr[0:nsecLen] + zero[nsecLen]
	} else {
		nsecStr = nsecStr[0:4]
	}
	timeStr := now.Format(base_format) + "." + nsecStr
	return timeStr
}
