# go-log [![GoDoc](https://godoc.org/github.com/golangly/log?status.svg)](http://godoc.org/github.com/golangly/log) [![Report card](https://goreportcard.com/badge/github.com/golangly/log)](https://goreportcard.com/report/github.com/golangly/log) [![Sourcegraph](https://sourcegraph.com/github.com/golangly/log/-/badge.svg)](https://sourcegraph.com/github.com/golangly/log?badge)

Go logging done right (well, you know...)

## Usage

```go
import "github.com/golangly/log"
...
log.Trace("Hello!")
log.Debug("Hello!")
log.Info("Hello!")
log.Warn("Hello!")
log.Error("Hello!")
log.Panic("Hello!")
log.Fatal("Hello!")
```

So far, so good - right? well that's basic usage identical to most other loggers. In fact, this logging implementation is probably even a bit slower than loggers such as zerolog or ruslog. There are, however, a few extra features that make this library both opinionated and more compatible to my view of logging:

* Accept recovered panic objects (via `recover()` in defer functions usually) and infer automatically whether they are errors or not, and treat them accordingly (it's still a panic, of course and treated as such also).
* Automatic configuration from environment variables since you want to configure logging before everything else, including before command line parsing.
* Only print to console. Use other tools to route logging to other places.
* Supports [`go-errors`](https://github.com/golangly/errors) for error stack traces and error causes automatically.
* Super simple.

## Configuration

All environment variable names & values are case insensitive.

Set `LOG_PRETTY` to `1`, `y`, `yes` or `true` to enable pretty print (by default `JSON` logging is used).

Set `LOG_LEVEL` to `disabled`, `panic` or `fatal` to disable logging (except panics & fatal events), or to `error`, `warn`, `info`, `debug` or `trace` to set the log level accordingly.
  
## Contributing

Please read the [Code of Conduct](.github/CODE_OF_CONDUCT.md) & [Contributing](.github/CONTRIBUTING.md) documents.

## License

[GNUv3](./LICENSE)
