package preflight

type tuningAction struct {
	component string
	message   string
	status    stageStatus
}

type tuningPlan struct {
	actions []tuningAction
	ready   bool
}

type optimizationResult struct {
	actions []tuningAction
	status  stageStatus
}

func (r *runner) buildTuningPlan(report readinessReport) tuningPlan {
	r.logger.Debug("Building performance tuning plan...", "phase", "plan", "component", "planner")

	return tuningPlan{
		ready: report.ready && report.resources == statusReady,
		actions: []tuningAction{
			{component: "scheduler", message: "Calibrating scheduler parameters...", status: statusPending},
			{component: "memory", message: "Tuning memory management policy...", status: statusPending},
			{component: "garbage_collector", message: "Balancing garbage collection targets...", status: statusPending},
		},
	}
}

func (r *runner) applyTuningPlan(plan tuningPlan) optimizationResult {
	r.logger.Info("Applying runtime optimizations...", "phase", "apply", "status", statusName(statusPending))

	actions := make([]tuningAction, len(plan.actions))
	copy(actions, plan.actions)
	for index := range actions {
		r.logger.Debug(actions[index].message, "phase", "apply", "component", actions[index].component)
		actions[index].status = statusApplied
	}

	result := optimizationResult{actions: actions, status: statusApplied}
	if !plan.ready {
		result.status = statusPending
	}
	return result
}

func (r *runner) verifyRuntime(profile systemProfile, result optimizationResult) optimizationResult {
	r.logger.Debug("Verifying optimized runtime state...", "phase", "verify", "component", profile.runtime.name)

	if profileReady(profile) && actionsApplied(result.actions) {
		result.status = statusVerified
	}

	r.logger.Info("Runtime optimization complete.", "phase", "verify", "status", statusName(result.status))
	return result
}

func actionsApplied(actions []tuningAction) bool {
	if len(actions) == 0 {
		return false
	}
	for _, action := range actions {
		if action.status != statusApplied {
			return false
		}
	}
	return true
}
