package app

import (
	"context"
	"log/slog"
	"sort"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/glacius-labs/schedulerd/internal/rules/filter"
	"github.com/glacius-labs/schedulerd/internal/rules/score"
)

type App struct {
	filterTree filter.Rule
	scoreTree  score.Rule
	logger     *slog.Logger
}

func NewApp(f filter.Rule, s score.Rule, l *slog.Logger) *App {
	return &App{
		filterTree: f,
		scoreTree:  s,
		logger:     l,
	}
}

func (a *App) Evaluate(ctx context.Context, workload domain.Workload, workers []domain.Worker) ([]domain.AssignmentResult, error) {
	var results []domain.AssignmentResult

	for _, w := range workers {
		ok, err := a.filterTree.Evaluate(ctx, workload, w)
		if err != nil {
			a.logger.Error("failed to evaluate filter tree", "error", err)
			continue
		}
		if !ok {
			continue
		}

		scoreVal, err := a.scoreTree.Evaluate(ctx, workload, w)
		if err != nil {
			a.logger.Error("failed to evaluate score tree", "error", err)
			continue
		}

		results = append(results, domain.AssignmentResult{
			WorkerID: w.ID,
			Score:    scoreVal,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}
