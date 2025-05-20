package rules

import (
	"bytes"
	"fmt"
	"os"

	"github.com/glacius-labs/schedulerd/internal/rules/filter"
	"github.com/glacius-labs/schedulerd/internal/rules/score"
	"github.com/google/cel-go/cel"
	"gopkg.in/yaml.v3"
)

func LoadRuleSet(filePath string) (*ParsedRules, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rule file: %w", err)
	}

	var node yaml.Node
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&node); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	var ruleSet RuleSet
	if err := node.Decode(&ruleSet); err != nil {
		return nil, fmt.Errorf("failed to extract rules: %w", err)
	}

	env, err := cel.NewEnv(
		cel.Variable("workload", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("worker", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	filterTree, err := filter.Parse(ruleSet.Filters, env)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	scoreTree, err := score.Parse(ruleSet.Scorers, env)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scorers: %w", err)
	}

	return &ParsedRules{
		FilterTree: filterTree,
		ScoreTree:  scoreTree,
	}, nil
}
