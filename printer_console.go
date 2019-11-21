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

	"github.com/huandu/go-tls"
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

func prettyPrint(c Logger, msg string) {
	needsQuote := func(s string) bool {
		for i := range s {
			if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
				return true
			}
		}
		return false
	}
	colorize := func(s interface{}, c int, disabled bool) string {
		if disabled {
			return fmt.Sprintf("%s", s)
		}
		return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
	}

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
	buf.WriteByte(' ')

	// Print process & thread IDs
	buf.WriteString(colorize(fmt.Sprintf(" pid=%d", pid), colorCyan, false))
	buf.WriteString(colorize(fmt.Sprintf(" tid=%d", tls.ID()), colorCyan, false))

	// Print context
	if len(c.context) > 0 {
		buf.WriteByte(' ')
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
			if needsQuote(stringValue) {
				buf.WriteString(strconv.Quote(stringValue))
			} else {
				buf.WriteString(stringValue)
			}
		}
	}
	buf.WriteByte('\n')

	// Print error
	if c.err != nil {
		_, _ = fmt.Fprintf(buf, "%+v\n", c.err)
	}

	// Write to output stream
	if _, err := buf.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
}
