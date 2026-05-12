package clog

import (
	"os"
	"strings"
	"time"

	"github.com/tobiashort/th-utils/lib/cfmt"
)

const (
	LevelDebug = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
)

func EmptyString() string {
	return ""
}

var (
	DebugString = func() string { return cfmt.Sprint("#B{DEBUG}") }
	InfoString  = func() string { return cfmt.Sprint("#bB{INFO}") }
	WarnString  = func() string { return cfmt.Sprint("#yB{WARN}") }
	ErrorString = func() string { return cfmt.Sprint("#rB{ERROR}") }
)

var (
	Level           = LevelInfo
	Output          = os.Stderr
	TimestampEnable = false
)

func keyValues(args ...any) string {
	sb := strings.Builder{}
	isVal := false
	for _, arg := range args {
		if isVal {
			cfmt.Fprintf(&sb, "%v ", arg)
			isVal = false
		} else {
			cfmt.Cfprintf(&sb, "c", "%s=", arg)
			isVal = true
		}
	}
	return sb.String()
}

func print(levelString string, msg string) {
	if levelString != "" {
		msg = cfmt.Sprintf("%s %s", levelString, msg)
	}
	if TimestampEnable {
		msg = cfmt.Sprintf("%s %s", time.Now().Format(time.RFC822), msg)
	}
	msg = strings.TrimRightFunc(msg, func(r rune) bool { return r == '\r' || r == '\n' })
	cfmt.Fprintln(Output, msg)
}

func Debug(args ...any) {
	if Level != LevelDebug {
		return
	}
	print(DebugString(), cfmt.Sprintln(args...))
}

func Debugf(format string, args ...any) {
	if Level != LevelDebug {
		return
	}
	print(DebugString(), cfmt.Sprintf(format, args...))
}

func Debugs(msg string, args ...any) {
	if Level != LevelDebug {
		return
	}
	print(DebugString(), cfmt.Sprintln(msg, keyValues(args...)))
}

func Info(args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(InfoString(), cfmt.Sprintln(args...))
}

func Infof(format string, args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(InfoString(), cfmt.Sprintf(format, args...))
}

func Infos(msg string, args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(InfoString(), cfmt.Sprintln(msg, keyValues(args...)))
}

func Warn(args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(WarnString(), cfmt.Sprintln(args...))
}

func Warnf(format string, args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(WarnString(), cfmt.Sprintf(format, args...))
}

func Warns(msg string, args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	print(WarnString(), cfmt.Sprintln(msg, keyValues(args...)))
}

func Error(args ...any) {
	print(ErrorString(), cfmt.Sprintln(args...))
}

func Errorf(format string, args ...any) {
	print(ErrorString(), cfmt.Sprintf(format, args...))
}

func Errors(msg string, args ...any) {
	print(ErrorString(), cfmt.Sprintln(msg, keyValues(args...)))
}
