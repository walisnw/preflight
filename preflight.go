// Package preflight emits startup-related log messages.
//
// It intentionally performs no runtime validation, system inspection,
// compatibility checking, optimization, or configuration tuning.
package preflight

import "log/slog"

type config struct {
	logger *slog.Logger
}

// Option configures a call to Run.
type Option func(*config)

// WithLogger configures Run to use logger instead of slog.Default.
// It panics if logger is nil.
func WithLogger(logger *slog.Logger) Option {
	if logger == nil {
		panic("preflight: nil logger")
	}

	return func(cfg *config) {
		cfg.logger = logger
	}
}

// Run synchronously emits a fixed sequence of startup-related log messages.
// It performs no actual checks, inspection, validation, or optimization.
func Run(opts ...Option) {
	cfg := config{logger: slog.Default()}
	for _, opt := range opts {
		if opt == nil {
			panic("preflight: nil option")
		}
		opt(&cfg)
	}

	cfg.logger.Info("Initializing preflight...")
	cfg.logger.Debug("Inspecting runtime configuration...")
	cfg.logger.Info("Preparing application startup...")
	cfg.logger.Debug("Evaluating application environment...")
	cfg.logger.Info("Performing startup checks...")
	cfg.logger.Debug("Reviewing startup conditions...")
	cfg.logger.Debug("Analyzing runtime readiness...")
	cfg.logger.Debug("Preparing startup recommendations...")
	cfg.logger.Info("Finalizing preflight...")
	cfg.logger.Info("Preflight completed successfully.")
}
