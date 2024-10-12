package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cloudokyo/env"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type (
	// An alias of zerolog.Event
	Event = zerolog.Event

	// An alias of zerolog.Logger
	Logger = zerolog.Logger
)

const (
	// A log field
	MSG = "msg"

	// A log field
	ANY = "any"
)

// A default instance of logger
var std Logger

var (
	// An alias of log.Print function
	Print = log.Print

	// An alias of log.Printf function
	Printf = log.Printf

	// An alias of log.Println function
	Println = log.Println
)

var (
	// Log level, default=DEBUG
	//  - UAT: DEBUG
	//  - PROD: INFO
	//
	// Env
	// 	LOG_LEVEL=DEBUG
	LogLevel = env.Get("LOG_LEVEL", "DEBUG")

	// Log trace mode.
	//  - If TRUE: write the message in INFO level
	//  - If FALSE: write the message in DEBUG level
	//
	// Env
	// 	LOG_TRACE=false
	LogTrace = env.GetBool("LOG_TRACE", false)

	// Log error stack trace.
	//  - If TRUE: append the stack trace into log message
	//  - If FALSE: write the logs as normal
	//
	// Env
	// 	LOG_STACK_TRACE=true
	LogStack = env.GetBool("LOG_STACK_TRACE", true)

	// Log to console mode.
	// - If TRUE: we write the logs to console with TEXT format, use on LOCALLY
	// - If FALSE: we write the logs to console in JSON format, use on remote
	//
	// Env
	// 	LOG_CONSOLE=false
	LogConsole = env.GetBool("LOG_CONSOLE", false)

	// The console log color
	//
	// Env
	// 	LOG_NOCOLOR=true
	LogNoColor = env.GetBool("LOG_NOCOLOR", true)
)

var (
	// The default log output is os.Stdout
	OutputDefault = os.Stdout

	// The default console log output
	OutputConsole = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    LogNoColor,
		TimeFormat: zerolog.TimeFieldFormat,
	}
)

func init() {
	// Trace the configs
	log.Println("- log.level:", LogLevel)
	log.Println("- log.trace:", LogTrace)
	log.Println("- log.stack:", LogStack)
	log.Println("- log.console:", LogConsole)
	log.Println("- log.nocolor:", LogNoColor)

	// Init the standard JSON logger instance
	level, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		log.Fatalf("%s --> zerolog.ParseLevel(%s)", err, LogLevel)
	}

	// Override zerolog configs
	zerolog.ErrorStackMarshaler = MarshalStack
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.999"
	zerolog.MessageFieldName = MSG

	var output io.Writer

	if LogConsole {
		output = OutputConsole
	} else {
		output = OutputDefault
	}

	// Init the standard log instance
	std = zerolog.New(output).
		Level(level).
		With().
		CallerWithSkipFrameCount(2).
		Timestamp().
		Logger()
}

// Get the default logger instance
//
// Example
//
//	log := log.Default()
func Default() *Logger {
	return &std
}

// Set the log output
//
// Example
//
//	log.Output(os.Stdout)
func Output(out io.Writer) Logger {
	std = std.Output(out)
	return std
}

// With creates a child logger with the field added to its context.
//
// Example
//
//	log.With()
func With() zerolog.Context {
	return std.With()
}

// Ctx returns the Logger associated with the ctx.
// If no logger is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) Logger {
	return std.Level(level)
}

// Write a log to output with TRACE level
//
// Example
//
//	log.Trace("hello", "world")
func Trace(v ...any) *Event {
	return Log(std.Trace(), v...)
}

// Write a log with format to output with TRACE level
//
// Example
//
//	log.Tracef("hello %s", "world")
func Tracef(format string, v ...any) {
	Logf(std.Trace(), format, v...)
}

// Write a log to output with DEBUG level
//
// Example
//
//	log.Debug("hello", "world").Send()
func Debug(v ...any) *Event {
	return Log(std.Debug(), v...)
}

// Write a log with format to output with DEBUG level
//
// Example
//
//	log.Debugf("hello %s", "world")
func Debugf(format string, v ...any) {
	Logf(std.Debug(), format, v...)
}

// Write a log to output with INFO level
//
// Example
//
//	log.Info("hello", "world").Send()
func Info(v ...any) *Event {
	return Log(std.Info(), v...)
}

// Write a log with format to output with INFO level
//
// Example
//
//	log.Infof("hello %s", "world")
func Infof(format string, v ...any) {
	Logf(std.Info(), format, v...)
}

// Write a log to output with WARN level
//
// Example
//
//	log.Warn("hello", "world").Send()
func Warn(v ...any) *Event {
	return Log(std.Warn(), v...)
}

// Write a log with format to output with WARN level
//
// Example
//
//	log.Warnf("hello %s", "world")
func Warnf(format string, v ...any) {
	Logf(std.Warn(), format, v...)
}

// Write a log to output with ERROR level
//
// Example
//
//	log.Error("hello", "world").Send()
func Error(v ...any) *Event {
	return Log(std.Error(), v...)
}

// Write a log with format to output with ERROR level
//
// Example
//
//	log.Errorf("hello %s", "world")
func Errorf(format string, v ...any) {
	Logf(std.Error(), format, v...)
}

// Write a log to output with PANIC level and PANIC
//
// Example
//
//	log.Panic("hello", "world").Send()
func Panic(v ...any) *Event {
	return Log(std.Panic().Stack(), v...)
}

// Write a log with format to output with PANIC level and PANIC
//
// Example
//
//	log.Panicf("hello %s", "world")
func Panicf(format string, v ...any) {
	Logf(std.Panic(), format, v...)
}

// Write a log to output with FATAL level and call os.Exit(1)
//
// Example
//
//	log.Fatal("hello", "world").Send()
func Fatal(v ...any) *Event {
	return Log(std.Fatal().Stack(), v...)
}

// Write a log with format to output with FATAL level and call os.Exit(1)
//
// Example
//
//	log.Fatalf("hello %s", "world")
func Fatalf(format string, v ...any) {
	Logf(std.Fatal(), format, v...)
}

// Write an event to log output
//
// Param
//
//	event: The zerolog event
//	args: The data injected to the log
//
// Example
//
//	log.Log(event, "Test log event message").Send()
func Log(event *Event, args ...any) *Event {
	msgs := []string{}
	for _, arg := range args {
		switch value := arg.(type) {
		case string:
			msgs = append(msgs, value)
		case error:
			// Check log stack mode
			if LogStack {
				event.Stack()
			}

			// Append error to log
			err := errors.WithStack(value)
			event.Err(err)
			event.AnErr(zerolog.ErrorFieldName, err)
		case context.Context:
			event.Ctx(value)
			if data := Detach(value); len(data) > 0 {
				event.Fields(data.Value())
			}
		case bool:
			if value {
				event.Stack()
			}
		case *bool:
			if *value {
				event.Stack()
			}
		default:
			event.Any(ANY, value)
		}
	}

	if len(msgs) > 0 {
		event.Str(MSG, strings.Join(msgs, " "))
	}

	return event
}

// Write an event to log output with format
//
// Param
//
//	event: The zerolog event
//	args: The data injected to the log
//
// Example
//
//	log.Logf(event, "Test log event message")
func Logf(event *Event, format string, args ...any) {
	values := []any{}

	// Check each arg is context
	for _, arg := range args {
		switch value := arg.(type) {
		case error:
			err := errors.Wrap(value, "")
			event.Err(err)
			event.AnErr(zerolog.ErrorFieldName, err)
			values = append(values, value)
		case context.Context:
			if data := Detach(value); data != nil {
				event.Fields(data.Value())
			}
		default:
			values = append(values, value)
		}
	}

	// Print the log
	event.CallerSkipFrame(2).Msgf(format, values...)
}

// Cast an object to string with JSON or formatter
//
// Example
//
//	log.String(data) --> string
func String(v any) string {
	if data, err := json.Marshal(v); err == nil {
		return string(data)
	}
	return fmt.Sprintf("%#v", v)
}
