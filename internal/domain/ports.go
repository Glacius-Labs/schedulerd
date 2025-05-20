package domain

type RuleEvaluator interface {
	EvaluateFilters(workload Workload, worker Worker) bool
	EvaluateScore(workload Workload, worker Worker) float64
}
