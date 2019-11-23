package log

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
)

type Level uint

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelDisabled
)

var (
	Root    = &Logger{context: make(map[string]interface{}, 0), level: LevelInfo}
	Printer func(c Logger, msg string)
	bufPool = sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 100)) }}
	pid     = os.Getpid()
	cwd     string
)

func init() {
	if _cwd, err := os.Getwd(); err != nil {
		cwd = _cwd
	}
	pretty := strings.ToLower(os.Getenv("LOG_PRETTY"))
	if pretty == "1" || pretty == "yes" || pretty == "true" || pretty == "y" {
		Printer = PrettyPrinter
	} else {
		Printer = JSONPrinter
	}

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "disabled", "panic", "fatal":
		Root = Root.WithLevel(LevelDisabled)
	case "error":
		Root = Root.WithLevel(LevelError)
	case "warn":
		Root = Root.WithLevel(LevelWarn)
	case "info", "":
		Root = Root.WithLevel(LevelInfo)
	case "debug":
		Root = Root.WithLevel(LevelDebug)
	case "trace":
		Root = Root.WithLevel(LevelTrace)
	default:
		Fatal("log level must be one of: disabled, panic, fatal, error, warn, info, debug or trace (case-insensitive)")
	}
}

func Writer() io.Writer                        { return Root.Writer() }
func With(key string, val interface{}) *Logger { return Root.With(key, val) }
func WithErr(err error) *Logger                { return Root.WithErr(err) }
func WithPanic(recovered interface{}) *Logger  { return Root.WithPanic(recovered) }
func WithLevel(level Level) *Logger            { return Root.WithLevel(level) }
func Trace(msg string)                         { Root.Trace(msg) }
func Debug(msg string)                         { Root.Debug(msg) }
func Info(msg string)                          { Root.Info(msg) }
func Warn(msg string)                          { Root.Warn(msg) }
func Error(msg string)                         { Root.Error(msg) }
func Panic(msg string)                         { Root.Panic(msg) }
func Fatal(msg string)                         { Root.Fatal(msg) }
