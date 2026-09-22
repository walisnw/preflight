// Package preflight emits a deterministic sequence of startup preflight log
// records. It does not inspect or modify the host system or Go runtime.
package preflight

import "log/slog"

// Run synchronously emits the preflight sequence through slog.Default.
// It does not perform system checks or runtime tuning.
func Run() {
	newRunner(slog.Default()).run()
}
