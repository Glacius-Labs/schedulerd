package score

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Weighted struct {
	Weight float64
	Rule   Rule
}

func (r *Weighted) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (float64, error) {
	val, err := r.Rule.Evaluate(ctx, w, wk)
	if err != nil {
		return 0.0, err
	}
	return val * r.Weight, nil
}

func ParseWeighted(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Weight float64   `yaml:"weight"`
		Rule   yaml.Node `yaml:"rule"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid weighted rule: %w", err)
	}

	child, err := Parse(raw.Rule, env)
	if err != nil {
		return nil, fmt.Errorf("failed to parse child rule in weighted: %w", err)
	}

	return &Weighted{Weight: raw.Weight, Rule: child}, nil
}
