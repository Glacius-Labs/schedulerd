package score

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Avg struct {
	Rules []Rule
}

func (r *Avg) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (float64, error) {
	if len(r.Rules) == 0 {
		return 0.0, nil
	}
	var sum float64
	for _, rule := range r.Rules {
		score, err := rule.Evaluate(ctx, w, wk)
		if err != nil {
			return 0.0, err
		}
		sum += score
	}
	return sum / float64(len(r.Rules)), nil
}

func ParseAvg(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Rules []yaml.Node `yaml:"rules"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid avg rule: %w", err)
	}

	var children []Rule
	for _, child := range raw.Rules {
		rule, err := Parse(child, env)
		if err != nil {
			return nil, fmt.Errorf("avg child parse error: %w", err)
		}
		children = append(children, rule)
	}
	return &Avg{Rules: children}, nil
}
