package errors

import (
	errors2 "errors"
	"fmt"
	"runtime"
)

func WithCtx(text string) error {
	skip := 0
	for {
		skip++
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			text = fmt.Sprintf("%s\n  at %s:%d", text, file, line)
		} else {
			break
		}
	}
	return errors2.New(text)
}

func WithCtxf(format string, a ...any) error {
	skip := 0
	for {
		skip++
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			format = fmt.Sprintf("%s\n  at %s:%d", format, file, line)
		} else {
			break
		}
	}
	return fmt.Errorf(format, a...)
}
