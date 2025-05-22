package domain

type Workload struct {
	ID     string         `json:"id" yaml:"id"`
	Labels map[string]any `json:"labels" yaml:"labels"`
}

type Worker struct {
	ID     string         `json:"id" yaml:"id"`
	Labels map[string]any `json:"labels" yaml:"labels"`
}

type AssignmentResult struct {
	WorkerID    string   `json:"worker_id" yaml:"worker_id"`
	Score       float64  `json:"score" yaml:"score"`
	Explanation []string `json:"explanation,omitempty" yaml:"explanation,omitempty"`
}
