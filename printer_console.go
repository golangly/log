package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//noinspection GoUnusedConst
const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func PrettyPrinter(c Logger, msg string) {

	// Checking a buffer from the pool
	var buf = bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	// Print timestamp
	buf.WriteString(colorize(time.Now().Format("15:04:05.000"), colorDarkGray, false))
	buf.WriteByte(' ')

	// Print level
	switch c.level {
	case LevelTrace:
		buf.WriteString(colorize("TRC", colorDarkGray, false))
	case LevelDebug:
		buf.WriteString(colorize("DBG", colorDarkGray, false))
	case LevelInfo:
		buf.WriteString(colorize("INF", colorGreen, false))
	case LevelWarn:
		buf.WriteString(colorize(colorize("WRN", colorYellow, false), colorBold, false))
	case LevelError:
		buf.WriteString(colorize(colorize("ERR", colorRed, false), colorBold, false))
	default:
		buf.WriteString(colorize(colorize("???", colorRed, false), colorBold, false))
	}
	buf.WriteByte(' ')

	// Print caller
	callerBlurb := "(unknown caller)"
	if _, file, line, ok := runtime.Caller(0); ok {
		caller := strings.TrimPrefix(strings.TrimPrefix(file, cwd), "/") + ":" + strconv.Itoa(line)
		callerBlurb = caller[len(caller)-15:]
	}
	buf.WriteString(fmt.Sprintf("%15.15s", colorize(callerBlurb[len(callerBlurb)-15:], colorBold, false)))
	buf.WriteString(colorize(" >", colorCyan, false))
	buf.WriteByte(' ')

	// Print message
	buf.WriteString(msg)
	buf.WriteString(" (")
	buf.WriteString(colorize(fmt.Sprintf("pid=%d", pid), colorCyan, false))

	// Print context
	if len(c.context) > 0 {
		for key, value := range c.context {
			buf.WriteString(colorize(fmt.Sprintf(" %s=", key), colorCyan, false))
			var stringValue string
			switch v := value.(type) {
			default:
				if b, err := json.Marshal(v); err != nil {
					stringValue = colorize(fmt.Sprintf("[error: %v]", err), colorRed, false)
				} else {
					stringValue = string(b)
				}
			}
			buf.WriteString(stringValue)
		}
	}
	buf.WriteString(")\n")

	// Print error
	if c.err != nil {
		// TODO: only use '%+.2v' on golangly errors; use '%v' otherwise
		_, _ = fmt.Fprintf(buf, "%+.2v\n", c.err)
	}

	// Write to output stream
	if _, err := buf.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}
