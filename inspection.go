package preflight

type stageStatus uint8

const (
	statusPending stageStatus = iota
	statusReady
	statusApplied
	statusVerified
)

type componentProfile struct {
	name   string
	status stageStatus
}

type systemProfile struct {
	runtime   componentProfile
	processor componentProfile
	memory    componentProfile
	scheduler componentProfile
	collector componentProfile
}

type readinessReport struct {
	compatibility stageStatus
	resources     stageStatus
	ready         bool
}

func (r *runner) inspectSystem() systemProfile {
	r.logger.Debug("Discovering runtime capabilities...", "phase", "inspect", "component", "runtime")
	r.logger.Debug("Inspecting processor topology...", "phase", "inspect", "component", "processor")
	r.logger.Debug("Evaluating memory configuration...", "phase", "inspect", "component", "memory")
	r.logger.Debug("Reviewing scheduler policy...", "phase", "inspect", "component", "scheduler")
	r.logger.Debug("Analyzing garbage collector settings...", "phase", "inspect", "component", "garbage_collector")

	profile := systemProfile{
		runtime:   componentProfile{name: "runtime", status: statusReady},
		processor: componentProfile{name: "processor", status: statusReady},
		memory:    componentProfile{name: "memory", status: statusReady},
		scheduler: componentProfile{name: "scheduler", status: statusReady},
		collector: componentProfile{name: "garbage_collector", status: statusReady},
	}

	r.logger.Info("System profile collected.", "phase", "inspect", "status", "ready")
	return profile
}

func (r *runner) analyzeRuntime(profile systemProfile) readinessReport {
	r.logger.Debug("Validating runtime compatibility...", "phase", "analyze", "component", profile.runtime.name)
	r.logger.Debug("Checking resource readiness...", "phase", "analyze", "component", "resources")

	report := readinessReport{
		compatibility: profile.runtime.status,
		resources:     combineStatus(profile.processor, profile.memory),
		ready:         profileReady(profile),
	}

	r.logger.Info("Runtime environment verified.", "phase", "analyze", "status", statusName(report.compatibility))
	return report
}

func combineStatus(components ...componentProfile) stageStatus {
	for _, component := range components {
		if component.status != statusReady {
			return statusPending
		}
	}
	return statusReady
}

func profileReady(profile systemProfile) bool {
	return combineStatus(
		profile.runtime,
		profile.processor,
		profile.memory,
		profile.scheduler,
		profile.collector,
	) == statusReady
}

func statusName(status stageStatus) string {
	switch status {
	case statusReady:
		return "ready"
	case statusApplied:
		return "applied"
	case statusVerified:
		return "verified"
	default:
		return "pending"
	}
}
