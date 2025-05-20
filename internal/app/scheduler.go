package app

import (
	"context"
	"sort"

	"github.com/glacius-labs/schedulerd/internal/domain"
	"github.com/glacius-labs/schedulerd/internal/rules/filter"
	"github.com/glacius-labs/schedulerd/internal/rules/score"
)

type Scheduler struct {
	FilterTree filter.Rule
	ScoreTree  score.Rule
}

func NewScheduler(f filter.Rule, s score.Rule) *Scheduler {
	return &Scheduler{
		FilterTree: f,
		ScoreTree:  s,
	}
}

func (s *Scheduler) Evaluate(ctx context.Context, workload domain.Workload, workers []domain.Worker) ([]domain.AssignmentResult, error) {
	var results []domain.AssignmentResult

	for _, w := range workers {
		ok, err := s.FilterTree.Evaluate(ctx, workload, w)
		if err != nil {
			// Optionally log and skip this worker
			continue
		}
		if !ok {
			continue
		}

		scoreVal, err := s.ScoreTree.Evaluate(ctx, workload, w)
		if err != nil {
			// Optionally log and skip
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
