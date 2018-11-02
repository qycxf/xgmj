//
// Author: leafsoar
// Date: 2016-06-24 15:18:11
//

// copy from name5566/leaf/log

package logs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

// 常用 Custom tag
const (
	PlayerTag  = "[player ] "
	AnalyseTag = "[analyse] "
	TableTag   = "[table  ] "
	NetworkTag = "[network] "
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
	warningLevel = 4
	customLevel  = 5
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printWarningLevel = "[warning] "
	printFatalLevel   = "[fatal  ] "
)

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File

	lock    sync.Mutex
	oldDate time.Time
	logfile string
	logpath string
}

func New(strLevel string, pathname string) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File
	if pathname != "" {
		now := time.Now()

		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())

		file, err := os.Create(path.Join(pathname, filename))
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", log.LstdFlags)
		baseFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.logpath = pathname
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	logger.oldDate = time.Now()

	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	// 分割日志逻辑
	if logger.baseFile != nil {
		now := time.Now()
		// 如果大于指定时间间隔，创建新文件
		if logger.oldDate.Add(time.Hour*6).Unix() < now.Unix() {
			logger.lock.Lock()
			defer logger.lock.Unlock()
			logger.oldDate = now
			logger.baseFile.Close()

			filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second())

			file, err := os.Create(path.Join(logger.logpath, filename))
			if err != nil {
				fmt.Println("创建日志文件失败")
			}

			logger.baseLogger = log.New(file, "", log.LstdFlags)
			logger.baseFile = file
			// 创建新文件的同时，删除旧有的 log
			logger.checkAndRemoveLogs(logger.logpath)
		}
	}

	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Printf(format, a...)

	if level == fatalLevel {
		// os.Exit(1)
	}
}

func (logger *Logger) checkAndRemoveLogs(path string) {
	// 保留今天和昨天的日志

}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func (logger *Logger) Warning(format string, a ...interface{}) {
	logger.doPrintf(warningLevel, printWarningLevel, format, a...)
}

func (logger *Logger) Custom(tag, format string, a ...interface{}) {
	logger.doPrintf(customLevel, tag, format, a...)
}

var gLogger, _ = New("debug", "")

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.Debug(format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.Release(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.Release(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.Fatal(format, a...)
}

func Warning(format string, a ...interface{}) {
	gLogger.Warning(format, a...)
}

func Custom(tag, format string, a ...interface{}) {
	gLogger.Custom(tag, format, a...)
}

func Close() {
	gLogger.Close()
}
