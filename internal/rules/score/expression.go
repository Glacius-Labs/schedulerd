package score

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

func (r *Expression) Evaluate(ctx context.Context, w domain.Workload, wk domain.Worker) (float64, error) {
	result, _, err := r.Program.Eval(map[string]any{
		"workload": w.Labels,
		"worker":   wk.Labels,
	})
	if err != nil {
		return 0.0, err
	}
	f, ok := result.Value().(float64)
	if !ok {
		return 0.0, fmt.Errorf("expression did not return a float64")
	}
	return f, nil
}

func ParseExpression(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Expression string `yaml:"expression"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid expression scorer: %w", err)
	}

	ast, iss := env.Compile(raw.Expression)
	if iss.Err() != nil {
		return nil, fmt.Errorf("CEL compile error: %w", iss.Err())
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("CEL program creation failed: %w", err)
	}

	return &Expression{Program: program}, nil
}
