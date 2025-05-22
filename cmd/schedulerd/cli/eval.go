package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/glacius-labs/schedulerd/internal/app"
	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/glacius-labs/schedulerd/internal/rules"
)

func EvalCmd(logger *slog.Logger) *cobra.Command {
	var configPath string
	var inputPath string

	cmd := &cobra.Command{
		Use:   "eval",
		Short: "Evaluate a workload against worker candidates using rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEval(configPath, inputPath, logger)
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "", "Path to rules.yaml")
	cmd.Flags().StringVar(&inputPath, "input", "", "Path to input.json")
	cmd.MarkFlagRequired("config")
	cmd.MarkFlagRequired("input")

	return cmd
}

func runEval(configPath, inputPath string, logger *slog.Logger) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	var input struct {
		Workload domain.Workload `json:"workload"`
		Workers  []domain.Worker `json:"workers"`
	}
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("invalid input JSON: %w", err)
	}

	ruleTrees, err := rules.LoadRuleSet(configPath)
	if err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}

	scheduler := app.NewApp(ruleTrees.FilterTree, ruleTrees.ScoreTree, logger)
	results, err := scheduler.Evaluate(context.Background(), input.Workload, input.Workers)
	if err != nil {
		return fmt.Errorf("evaluation failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
