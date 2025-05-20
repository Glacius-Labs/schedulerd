package filter

import (
	"context"

	"github.com/glacius-labs/schedulerd/internal/domain"
)

type Rule interface {
	Evaluate(ctx context.Context, workload domain.Workload, worker domain.Worker) (bool, error)
}
