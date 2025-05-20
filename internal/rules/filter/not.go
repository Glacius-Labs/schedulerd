package filter

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Not struct {
	Child Rule
}

func (r *Not) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (bool, error) {
	res, err := r.Child.Evaluate(ctx, w, wk)
	if err != nil {
		return false, err
	}
	return !res, nil
}

func ParseNot(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Rule yaml.Node `yaml:"rule"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid not rule: %w", err)
	}

	child, err := Parse(raw.Rule, env)
	if err != nil {
		return nil, fmt.Errorf("not rule child parse failed: %w", err)
	}
	return &Not{Child: child}, nil
}
