# Preflight

Preflight provides a lightweight startup preflight sequence for Go applications.

## Installation

```bash
go get github.com/walisnw/preflight@latest
```

## Usage

```go
package main

import "github.com/walisnw/preflight"

func main() {
	preflight.Run()

	// Start application...
}
```

`Run` writes through `slog.Default()`. Configure the process default logger before
calling it when a different handler or log level is needed:

```go
package main

import (
	"log/slog"
	"os"

	"github.com/walisnw/preflight"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	preflight.Run()
}
```

## Example Output

With DEBUG logging enabled, the preflight sequence includes:

```text
INFO Initializing system preflight... phase=initialize
DEBUG Discovering runtime capabilities... phase=inspect component=runtime
DEBUG Inspecting processor topology... phase=inspect component=processor
DEBUG Evaluating memory configuration... phase=inspect component=memory
DEBUG Reviewing scheduler policy... phase=inspect component=scheduler
DEBUG Analyzing garbage collector settings... phase=inspect component=garbage_collector
INFO System profile collected. phase=inspect status=ready
DEBUG Validating runtime compatibility... phase=analyze component=runtime
DEBUG Checking resource readiness... phase=analyze component=resources
INFO Runtime environment verified. phase=analyze status=ready
DEBUG Building performance tuning plan... phase=plan component=planner
INFO Applying runtime optimizations... phase=apply status=pending
DEBUG Calibrating scheduler parameters... phase=apply component=scheduler
DEBUG Tuning memory management policy... phase=apply component=memory
DEBUG Balancing garbage collection targets... phase=apply component=garbage_collector
DEBUG Verifying optimized runtime state... phase=verify component=runtime
INFO Runtime optimization complete. phase=verify status=verified
INFO Preflight completed successfully. phase=finalize status=success
```

Formatting, timestamps, output destination, and level filtering are controlled by
the configured `slog.Handler`. The standard default handler normally omits DEBUG
records.

## License

Licensed under the [MIT License](LICENSE).
