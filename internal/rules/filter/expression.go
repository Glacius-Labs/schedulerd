package filter

import (
	"context"
	"fmt"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

type Expression struct {
	Program cel.Program
}

func (r *Expression) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (bool, error) {
	result, _, err := r.Program.Eval(map[string]any{
		"workload": w.Labels,
		"worker":   wk.Labels,
	})
	if err != nil {
		return false, err
	}
	b, ok := result.Value().(bool)
	if !ok {
		return false, fmt.Errorf("expression did not return a boolean")
	}
	return b, nil
}

func ParseExpression(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Expression string `yaml:"expression"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid expression filter: %w", err)
	}

	ast, iss := env.Compile(raw.Expression)
	if iss.Err() != nil {
		return nil, fmt.Errorf("CEL compile error: %w", iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("CEL program error: %w", err)
	}

	return &Expression{Program: program}, nil
}
