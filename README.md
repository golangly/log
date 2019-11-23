# go-log [![GoDoc](https://godoc.org/github.com/arikkfir/go-log?status.svg)](http://godoc.org/github.com/arikkfir/go-log) [![Report card](https://goreportcard.com/badge/github.com/arikkfir/go-log)](https://goreportcard.com/report/github.com/arikkfir/go-log) [![Sourcegraph](https://sourcegraph.com/github.com/arikkfir/go-log/-/badge.svg)](https://sourcegraph.com/github.com/arikkfir/go-log?badge)

Go logging done right (well, you know...)

## Usage

```go
import "github.com/arikkfir/go-log"
...
log.Trace("Hello!")
log.Debug("Hello!")
log.Info("Hello!")
log.Warn("Hello!")
log.Error("Hello!")
log.Panic("Hello!")
log.Fatal("Hello!")
```

So far, so good - right? well that's basic usage identical to most other loggers. In factm, this logging implementation is probably even a bit slower than loggers such as zerolog or ruslog. There are, however, a few extra features that make this library both opinionated and more compatible to my view of logging:

* Accept recovered panic objects (via `recover()` in defer functions usually) and infer automatically whether they are errors or not, and treat them accordingly (it's still a panic, of course and treated as such also).
* Automatic configuration from environment variables since you want to configure logging before everything else, including before command line parsing.
* Only print to console. Use other tools to route logging to other places.
* Supports [`go-errors`](https://github.com/arikkfir/go-errors) for error stack traces and error causes automatically.
* Super simple.

## Configuration

Use `LOG_PRETTY` with values of `1`, `y`, `yes`, `true` (case insensitively) to pretty print (by default `JSON` logging is used).
  
## Contributing

Please read the [Code of Conduct](./docs/CODE_OF_CONDUCT.md) & [Contributing](./docs/CONTRIBUTING.md) documents.

## License

[GNUv3](./LICENSE)
