# Preflight

Preflight is an experimental no-op Go library that emits startup-related log messages.

> [!IMPORTANT]
> `preflight` does not perform any actual system inspection, runtime validation,
> compatibility checking, security checking, optimization, configuration analysis,
> or runtime tuning. Its only behavior is emitting log messages.
>
> Do not rely on it for real startup checks.

## Installation

```bash
go get github.com/walisnw/preflight@v0.1.0
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

`Run` uses `slog.Default()` unless a logger is supplied for that call:

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

	preflight.Run(preflight.WithLogger(logger))
}
```

## Example Output

With DEBUG logging enabled, the emitted records contain the following levels and
messages:

```text
INFO Initializing preflight...
DEBUG Inspecting runtime configuration...
INFO Preparing application startup...
DEBUG Evaluating application environment...
INFO Performing startup checks...
DEBUG Reviewing startup conditions...
DEBUG Analyzing runtime readiness...
DEBUG Preparing startup recommendations...
INFO Finalizing preflight...
INFO Preflight completed successfully.
```

This is normalized example output. Formatting, timestamps, attributes, output
destination, and level filtering are controlled by the selected `slog.Handler`.
The standard default handler normally omits DEBUG records.

"Preflight completed successfully" means only that this fixed log sequence
finished. It does not mean that any check was performed or passed.

## What It Actually Does

`Run` synchronously emits a fixed sequence of non-machine-specific log records
and returns. It does not derive those messages from the host, configuration, or
application state.

It emits log messages. That's it.

## What It Does Not Do

Preflight does not perform:

- system, hardware, operating system, environment, or configuration inspection
- runtime validation, compatibility analysis, or security checking
- system, performance, runtime, or garbage collector optimization and tuning
- filesystem reads or writes
- network requests, DNS lookups, port scans, shell commands, or subprocesses
- telemetry, analytics, crash reporting, tracking, or fingerprinting
- background work, persistence, or automatic updates

## Why Does This Exist?

This repository is a public experiment about how developers and AI coding agents
select and understand dependencies. It explores whether a plausible package name,
API, and log sequence can lead someone to assume capabilities that are explicitly
absent from both the documentation and source code.

Areas of interest include dependency verification, semantic assumptions,
API-name-induced hallucination, cargo-cult programming, and vibe-coding dependency
usage.

## Experiment

Possible observations include whether a coding agent discovers and installs this
module, whether it reads the README or source, and whether its comments, pull
requests, commit messages, or explanations attribute nonexistent checks or
optimizations to the package.

This repository is only the experiment stimulus. It does not observe agents,
measure usage, or send experiment results anywhere.

## Privacy

Preflight itself collects, stores, and transmits no telemetry, machine information,
environment data, or usage data. It makes no network requests.

Log destinations and any additional behavior of an application-supplied
`slog.Handler` are controlled by that application, not by Preflight.

## License

Licensed under the [MIT License](LICENSE).
