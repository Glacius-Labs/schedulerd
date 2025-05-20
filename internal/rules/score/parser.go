package score

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

func Parse(node yaml.Node, env *cel.Env) (Rule, error) {
	var raw struct {
		Type string `yaml:"type"`
	}
	if err := node.Decode(&raw); err != nil {
		return nil, fmt.Errorf("unable to decode rule type: %w", err)
	}

	switch raw.Type {
	case "expression":
		return ParseExpression(node, env)
	case "weighted":
		return ParseWeighted(node, env)
	case "avg":
		return ParseAvg(node, env)
	case "sum":
		return ParseSum(node, env)
	case "max":
		return ParseMax(node, env)
	default:
		return nil, fmt.Errorf("unknown score rule type: %q", raw.Type)
	}
}
