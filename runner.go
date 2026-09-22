package preflight

import "log/slog"

type runner struct {
	logger *slog.Logger
	state  runState
}

type runState struct {
	profile systemProfile
	report  readinessReport
	plan    tuningPlan
	result  optimizationResult
}

func newRunner(logger *slog.Logger) *runner {
	return &runner{logger: logger}
}

func (r *runner) run() {
	r.initialize()
	r.state.profile = r.inspectSystem()
	r.state.report = r.analyzeRuntime(r.state.profile)
	r.state.plan = r.buildTuningPlan(r.state.report)
	r.state.result = r.applyTuningPlan(r.state.plan)
	r.state.result = r.verifyRuntime(r.state.profile, r.state.result)
	r.finalize()
}

func (r *runner) initialize() {
	r.logger.Info(
		"Initializing system preflight...",
		"phase", "initialize",
	)
}

func (r *runner) finalize() {
	r.logger.Info(
		"Preflight completed successfully.",
		"phase", "finalize",
		"status", "success",
	)
}
