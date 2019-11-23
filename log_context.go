package log

import (
	"fmt"
	"io"
	"os"

	"github.com/arikkfir/go-errors"
)

type Logger struct {
	context map[string]interface{}
	err     error
	level   Level
}

type logWriter struct {
	logger *Logger
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
	}
	w.logger.Info(string(p))
	return
}

func (c *Logger) Writer() io.Writer { return &logWriter{logger: c} }

func (c *Logger) With(key string, val interface{}) *Logger {
	nc := &Logger{context: make(map[string]interface{}, len(c.context)), err: c.err, level: c.level}
	for k, v := range c.context {
		nc.context[k] = v
	}
	nc.context[key] = val
	return nc
}

func (c *Logger) WithErr(err error) *Logger {
	nc := &Logger{context: make(map[string]interface{}, len(c.context)), err: err, level: c.level}
	for k, v := range c.context {
		nc.context[k] = v
	}
	return nc
}

func (c *Logger) WithPanic(recovered interface{}) *Logger {
	var e errors.ErrorExt
	if err, ok := recovered.(error); ok {
		e = errors.Wrap(err, fmt.Sprintf("recovered panic: %v", recovered))
	} else {
		e = errors.New(fmt.Sprintf("recovered panic: %v", recovered))
	}
	e = e.AddTag("recovered", recovered)
	return c.WithErr(e)
}

func (c *Logger) WithLevel(level Level) *Logger {
	nc := &Logger{context: make(map[string]interface{}, len(c.context)), err: c.err, level: level}
	for k, v := range c.context {
		nc.context[k] = v
	}
	return nc
}

func (c *Logger) Trace(msg string) {
	if c.level <= LevelTrace {
		c.print(msg)
	}
}
func (c *Logger) Debug(msg string) {
	if c.level <= LevelDebug {
		c.print(msg)
	}
}
func (c *Logger) Info(msg string) {
	if c.level <= LevelInfo {
		c.print(msg)
	}
}
func (c *Logger) Warn(msg string) {
	if c.level <= LevelWarn {
		c.print(msg)
	}
}
func (c *Logger) Error(msg string) {
	if c.level <= LevelError {
		c.print(msg)
	}
}

func (c *Logger) Panic(msg string) { c.print(msg); panic(c) }

func (c *Logger) Fatal(msg string) { c.print(msg); os.Exit(1) }

func (c *Logger) print(msg string) { printer(*c, msg) }
