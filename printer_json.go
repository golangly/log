package log

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/huandu/go-tls"
)

func jsonPrint(c Logger, msg string) {

	// Create log entry map
	o := map[string]interface{}{
		"time": time.Now().Format("15:04:05.000"),
	}

	// Print timestamp
	o["time"] = time.Now().Format("15:04:05.000")

	// Print level
	switch c.level {
	case LevelTrace:
		o["level"] = "TRC"
	case LevelDebug:
		o["level"] = "DBG"
	case LevelInfo:
		o["level"] = "INF"
	case LevelWarn:
		o["level"] = "WRN"
	case LevelError:
		o["level"] = "ERR"
	default:
		o["level"] = "???"
	}

	// Print caller
	callerBlurb := "(unknown caller)"
	if _, file, line, ok := runtime.Caller(0); ok {
		caller := strings.TrimPrefix(strings.TrimPrefix(file, cwd), "/") + ":" + strconv.Itoa(line)
		callerBlurb = caller[len(caller)-15:]
	}
	o["caller"] = callerBlurb

	// Print message
	o["msg"] = msg

	// Print process & thread IDs
	o["pid"] = pid
	o["tid"] = tls.ID()

	// Print context
	if len(c.context) > 0 {
		for key, value := range c.context {
			o[key] = value
		}
	}

	// Print error
	if c.err != nil {
		o["error"] = fmt.Sprintf("%+v", c.err)
	}

	// Marshall to JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(o); err != nil {
		panic(err)
	}
}
