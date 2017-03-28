package logs

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var logChan chan string

func init() {
	logChan = make(chan string, 1024)

	go func() {
		for {
			<-logChan
			// 从Channel中取出来写文件

		}
	}()
}

func Debug(format string, args ...interface{}) {
	fun, file, line, ok := runtime.Caller(1)
	if ok {
		file := strings.TrimSuffix(path.Base(file), path.Ext(file))
		funcName := runtime.FuncForPC(fun).Name()

		t := time.Now().Format(time.Stamp)
		brush("["+t+"]", 58)
		brush(" "+"DEBUG"+" ", 58)
		brush(" "+file+" ", 58)
		brush(" "+funcName+"("+strconv.Itoa(line)+")|", 58)
		brush(fmt.Sprintf(format, args...), 10)
		fmt.Printf("\n")

		txts := fmt.Sprintf(format, args...)
		logChan <- formatLog(fatalLog, line, funcName, file, txts)
	}
}

func Info(format string, args ...interface{}) {
	fun, file, line, ok := runtime.Caller(1)
	if ok {
		file := strings.TrimSuffix(path.Base(file), path.Ext(file))
		funcName := runtime.FuncForPC(fun).Name()

		t := time.Now().Format(time.Stamp)
		brush("["+t+"]", 55)
		brush(" "+"INFO "+" ", 55)
		brush(" "+file+" ", 55)
		brush(" "+funcName+"("+strconv.Itoa(line)+")|", 55)
		brush(fmt.Sprintf(format, args...), 7)
		fmt.Printf("\n")

		txts := fmt.Sprintf(format, args...)
		logChan <- formatLog(fatalLog, line, funcName, file, txts)
	}
}

func Warn(format string, args ...interface{}) {
	fun, file, line, ok := runtime.Caller(1)
	if ok {
		file := strings.TrimSuffix(path.Base(file), path.Ext(file))
		funcName := runtime.FuncForPC(fun).Name()

		t := time.Now().Format(time.Stamp)
		brush("["+t+"]", 62)
		brush(" "+"WARN "+" ", 62)
		brush(" "+file+" ", 62)
		brush(" "+funcName+"("+strconv.Itoa(line)+")|", 62)
		brush(fmt.Sprintf(format, args...), 14)
		fmt.Printf("\n")

		txts := fmt.Sprintf(format, args...)
		logChan <- formatLog(fatalLog, line, funcName, file, txts)
	}
}

func Error(format string, args ...interface{}) {
	fun, file, line, ok := runtime.Caller(1)
	if ok {
		file := strings.TrimSuffix(path.Base(file), path.Ext(file))
		funcName := runtime.FuncForPC(fun).Name()

		t := time.Now().Format(time.Stamp)
		brush("["+t+"]", 61)
		brush(" "+"ERROR"+" ", 61)
		brush(" "+file+" ", 61)
		brush(" "+funcName+"("+strconv.Itoa(line)+")|", 61)
		brush(fmt.Sprintf(format, args...), 13)
		fmt.Printf("\n")

		txts := fmt.Sprintf(format, args...)
		logChan <- formatLog(fatalLog, line, funcName, file, txts)
	}
}

func Fatal(format string, args ...interface{}) {
	fun, file, line, ok := runtime.Caller(1)
	if ok {
		file := strings.TrimSuffix(path.Base(file), path.Ext(file))
		funcName := runtime.FuncForPC(fun).Name()

		t := time.Now().Format(time.Stamp)
		brush("["+t+"]", 76)
		brush(" "+"FATAL"+" ", 76)
		brush(" "+file+" ", 76)
		brush(" "+funcName+"("+strconv.Itoa(line)+")|", 76)
		brush(fmt.Sprintf(format, args...), 12)
		fmt.Printf("\n")

		txts := fmt.Sprintf(format, args...)
		logChan <- formatLog(fatalLog, line, funcName, file, txts)
	}
}

func formatLog(level severity, line int, funcName, file, texts string) string {

	t := time.Now().Format(time.Stamp)
	ss := fmt.Sprintf("[%s] %s [%s] %s(%d) | %s", t, severityName[level.get()], file,
		funcName, line, texts)
	return ss
}
