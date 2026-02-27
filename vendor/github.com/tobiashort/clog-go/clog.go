package clog

import (
	"os"
	"strings"

	"github.com/tobiashort/cfmt-go"
)

const (
	LevelDebug = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
)

var (
	DebugString = func() string { return cfmt.Sprint("#B{DEBUG}") }
	InfoString  = func() string { return cfmt.Sprint("#bB{INFO}") }
	WarnString  = func() string { return cfmt.Sprint("#yB{WARN}") }
	ErrorString = func() string { return cfmt.Sprint("#rB{ERROR}") }
)

var (
	Level  = LevelInfo
	Output = os.Stderr
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

func Debug(args ...any) {
	if Level != LevelDebug {
		return
	}
	cfmt.Fprintf(Output, "%s %s", DebugString(), cfmt.Sprintln(args...))
}

func Debugf(format string, args ...any) {
	if Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, DebugString(), cfmt.Sprintf(format, args...))
}

func Debugs(msg string, args ...any) {
	if Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, DebugString(), msg, keyValues(args...))
}

func Info(args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintf(Output, "%s %s", InfoString(), cfmt.Sprintln(args...))
}

func Infof(format string, args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, InfoString(), cfmt.Sprintf(format, args...))
}

func Infos(msg string, args ...any) {
	if Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, InfoString(), msg, keyValues(args...))
}

func Warn(args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintf(Output, "%s %s", WarnString(), cfmt.Sprintln(args...))
}

func Warnf(format string, args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, WarnString(), cfmt.Sprintf(format, args...))
}

func Warns(msg string, args ...any) {
	if Level != LevelWarn && Level != LevelInfo && Level != LevelDebug {
		return
	}
	cfmt.Fprintln(Output, WarnString(), msg, keyValues(args...))
}

func Error(args ...any) {
	cfmt.Fprintf(Output, "%s %s", ErrorString(), cfmt.Sprintln(args...))
}

func Errorf(format string, args ...any) {
	cfmt.Fprintln(Output, ErrorString(), cfmt.Sprintf(format, args...))
}

func Errors(msg string, args ...any) {
	cfmt.Fprintln(Output, ErrorString(), msg, keyValues(args...))
}
