package filter

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Or struct {
	Rules []Rule
}

func (r *Or) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (bool, error) {
	for _, rule := range r.Rules {
		res, err := rule.Evaluate(ctx, w, wk)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func ParseOr(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Rules []yaml.Node `yaml:"rules"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid or rule: %w", err)
	}

	children := make([]Rule, 0, len(raw.Rules))
	for _, child := range raw.Rules {
		r, err := Parse(child, env)
		if err != nil {
			return nil, fmt.Errorf("or rule child failed: %w", err)
		}
		children = append(children, r)
	}
	return &Or{Rules: children}, nil
}
