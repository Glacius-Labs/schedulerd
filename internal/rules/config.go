package rules

import (
	"github.com/glacius-labs/schedulerd/internal/rules/filter"
	"github.com/glacius-labs/schedulerd/internal/rules/score"
	"gopkg.in/yaml.v3"
)

type RuleSet struct {
	Filters yaml.Node `yaml:"filters"`
	Scorers yaml.Node `yaml:"scorers"`
}

type ParsedRules struct {
	FilterTree filter.Rule
	ScoreTree  score.Rule
}
