package filter

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

func Parse(node yaml.Node, env *cel.Env) (Rule, error) {
	var typeWrapper struct {
		Type string `yaml:"type"`
	}
	if err := node.Decode(&typeWrapper); err != nil {
		return nil, fmt.Errorf("failed to extract type: %w", err)
	}

	switch typeWrapper.Type {
	case "expression":
		return ParseExpression(node, env)
	case "and":
		return ParseAnd(node, env)
	case "or":
		return ParseOr(node, env)
	case "not":
		return ParseNot(node, env)
	default:
		return nil, fmt.Errorf("unsupported filter type: %s", typeWrapper.Type)
	}
}
