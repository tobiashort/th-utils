package assert

import (
	"fmt"
)

var PanicFunc = func(err error) {
	panic(err)
}

func Nil(v any, format string, args ...any) {
	if v != nil {
		PanicFunc(fmt.Errorf(format, args...))
	}
}

func NotNil(v any, format string, args ...any) {
	if v == nil {
		PanicFunc(fmt.Errorf(format, args...))
	}
}

func True(cond bool, format string, args ...any) {
	if !cond {
		PanicFunc(fmt.Errorf(format, args...))
	}
}

func False(cond bool, format string, args ...any) {
	if cond {
		PanicFunc(fmt.Errorf(format, args...))
	}
}
