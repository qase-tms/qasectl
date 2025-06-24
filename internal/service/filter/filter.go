package filter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qase-tms/qasectl/internal/models/plan"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	GetPlan(ctx context.Context, projectCode string, planID int64) (plan.PlanDetailed, error)
}

type Service struct {
	client client
}

// NewService creates a new Service
func NewService(client client) *Service {
	return &Service{client: client}
}

// GetFilteredResults returns the filtered results for the given plan ID and framework
func (s *Service) GetFilteredResults(ctx context.Context, project string, planID int64, framework string) (string, error) {
	const op = "result.service.getfilteredresults"
	logger := slog.With("op", op)

	logger.Debug("getting filtered results", "project", project, "planID", planID, "framework", framework)

	plan, err := s.client.GetPlan(ctx, project, planID)
	if err != nil {
		return "", err
	}

	if len(plan.Cases) == 0 {
		return "", fmt.Errorf("no cases found in plan")
	}

	switch framework {
	case "playwright":
		return prepareForPlaywright(plan.Cases), nil
	default:
		return "", fmt.Errorf("unsupported framework: %s", framework)
	}
}

// prepareForPlaywright prepares the IDs for playwright. Create regex pattern: "(Qase ID: 1|2|3|...)"
func prepareForPlaywright(IDs []int64) string {
	pattern := "(Qase ID: "
	for i, caseID := range IDs {
		if i > 0 {
			pattern += "|"
		}
		pattern += fmt.Sprintf("%d", caseID)
	}
	pattern += ")"

	return pattern
}
