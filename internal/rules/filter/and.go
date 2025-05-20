package filter

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type And struct {
	Rules []Rule
}

func (r *And) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (bool, error) {
	for _, rule := range r.Rules {
		res, err := rule.Evaluate(ctx, w, wk)
		if err != nil || !res {
			return false, err
		}
	}
	return true, nil
}

func ParseAnd(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Rules []yaml.Node `yaml:"rules"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid and rule: %w", err)
	}

	children := make([]Rule, 0, len(raw.Rules))
	for _, child := range raw.Rules {
		r, err := Parse(child, env)
		if err != nil {
			return nil, fmt.Errorf("and rule child failed: %w", err)
		}
		children = append(children, r)
	}
	return &And{Rules: children}, nil
}
