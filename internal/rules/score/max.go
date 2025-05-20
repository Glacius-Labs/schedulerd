package score

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Max struct {
	Rules []Rule
}

func (r *Max) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (float64, error) {
	if len(r.Rules) == 0 {
		return 0.0, nil
	}
	var max float64
	first := true
	for _, rule := range r.Rules {
		score, err := rule.Evaluate(ctx, w, wk)
		if err != nil {
			return 0.0, err
		}
		if first || score > max {
			max = score
			first = false
		}
	}
	return max, nil
}

func ParseMax(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Rules []yaml.Node `yaml:"rules"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid max rule: %w", err)
	}

	var children []Rule
	for _, child := range raw.Rules {
		rule, err := Parse(child, env)
		if err != nil {
			return nil, fmt.Errorf("max child parse error: %w", err)
		}
		children = append(children, rule)
	}
	return &Max{Rules: children}, nil
}
