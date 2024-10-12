package log

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return true
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

// MarshalStack implements pkg/errors stack trace marshaling.
//
// zerolog.ErrorStackMarshaler = MarshalStack
func MarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var sterr stackTracer
	var ok bool
	for err != nil {
		sterr, ok = err.(stackTracer)
		if ok {
			break
		}

		u, ok := err.(interface {
			Unwrap() error
		})
		if !ok {
			return nil
		}

		err = u.Unwrap()
	}
	if sterr == nil {
		return nil
	}

	// Get stack trace frame
	st := sterr.StackTrace()

	// Skip the last frame
	if c := len(st); c > 0 {
		st = st[:c-1]
	}

	// Init the state
	s := &state{}

	// Init the output
	out := make([]string, 0, len(st))

	// Process each frame to build the output
	for _, frame := range st {
		file := frameField(frame, s, 's')
		line := frameField(frame, s, 'd')
		fn := funcname(frameField(frame, s, 'n'))
		pkg := pkgname(file)

		// Skip log package frames
		if pkg == "log" {
			continue
		}

		// Skip other packages frames
		if strings.Index(file, "net/http") > 0 || strings.Index(file, "labstack/echo") > 0 {
			continue
		}

		out = append(out, fmt.Sprintf("%s.%s:%s:%s", pkg, fn, file, line))
	}
	return out
}

// Get the function name from the input
//
// Example
//
//	funcname("Log") --> Log
//	funcname("(*service).processEventResult") --> processEventResult
func funcname(name string) string {
	chunks := strings.Split(name, ").")
	if count := len(chunks); count > 0 {
		return chunks[count-1]
	} else {
		return name
	}
}

// Get the package name from a file path
//
// Example
//
//	pkgname("") --> main
//	pkgname("test.go") --> main
//	pkgname("pkg/test/log") --> log
//	pkgname("/go/pkg/mod/github.com/labstack/echo/v4@v4.12.0/echo.go") --> echo
func pkgname(filePath string) string {
	chunks := strings.Split(filePath, "/")
	count := len(chunks)
	if count > 1 {
		for i := count - 2; i >= 0; i-- {
			if val := chunks[i]; val == "" {
				continue
			} else if strings.Contains(val, ".") {
				continue // ignore version, e.g: v4@v4.12.0
			} else {
				return val
			}
		}

		return "main"
	} else {
		return strings.TrimSuffix(filePath, ".go")
	}
}
