package log

import (
	"fmt"
	"io"
	std "log"
	"os"
)

type LogColor int32

const (
	LOGCOLOR_RED    LogColor = 12
	LOGCOLOR_GREEN           = 10
	LOGCOLOR_YELLOW          = 14
	LOGCOLOR_GRAY            = 7
	LOGCOLOR_PINK            = 13
	LOGCOLOR_BLUE            = 9
	LOGCOLOR_WHITE           = 15
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

func New(out io.Writer, prefix string, flag int) *std.Logger {
	return std.New(out, prefix, flag)
}

var stdLog = New(os.Stderr, "", std.LstdFlags)

func SetOutput(w io.Writer) {
	stdLog.SetOutput(w)
}

func Flags() int {
	return stdLog.Flags()
}

func SetFlags(flag int) {
	stdLog.SetFlags(flag)
}

func Prefix() string {
	return stdLog.Prefix()
}

func SetPrefix(prefix string) {
	stdLog.SetPrefix(prefix)
}

func Print(v ...interface{}) {
	stdLog.Output(2, fmt.Sprint(v...))
}

func Debug(format string, v ...interface{}) {
	ColorOutput(LOGCOLOR_GRAY, 2, fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	ColorOutput(LOGCOLOR_GREEN, 2, fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	ColorOutput(LOGCOLOR_YELLOW, 2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	ColorOutput(LOGCOLOR_PINK, 2, fmt.Sprintf(format, v...))

}

func Fatal(format string, v ...interface{}) {
	ColorOutput(LOGCOLOR_RED, 2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func ColorOutput(color LogColor, calldepth int, s string) {
	stdLog.SetFlags(std.Lshortfile | std.Ldate | std.Ltime)
	stdLog.Output(calldepth+1, s)
}
