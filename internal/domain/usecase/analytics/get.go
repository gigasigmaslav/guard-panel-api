package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
)

type gaAnalyticsRepo interface {
	GetTaskDashboardKPI(
		ctx context.Context,
		from, to time.Time,
	) (contract.TaskDashboardKPI, error)

	GetTopExecutorByCompletedInPeriod(
		ctx context.Context,
		from, to time.Time,
	) (contract.TopExecutorByCompleted, error)
}

type GetAnalyticsUseCase struct {
	analyticsRepo gaAnalyticsRepo
}

func NewGetAnalyticsUseCase(
	analyticsRepo gaAnalyticsRepo,
) *GetAnalyticsUseCase {
	return &GetAnalyticsUseCase{
		analyticsRepo: analyticsRepo,
	}
}

func (ga *GetAnalyticsUseCase) Get(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (contract.TaskDashboardKPI, contract.TopExecutorByCompleted, error) {
	from, to, err := normalizeAnalyticsPeriod(from, to)
	if err != nil {
		return contract.TaskDashboardKPI{}, contract.TopExecutorByCompleted{}, err
	}

	kpi, err := ga.analyticsRepo.GetTaskDashboardKPI(ctx, from, to)
	if err != nil {
		return contract.TaskDashboardKPI{}, contract.TopExecutorByCompleted{},
			fmt.Errorf("get task dashboard kpi: %w", err)
	}

	top, err := ga.analyticsRepo.GetTopExecutorByCompletedInPeriod(ctx, from, to)
	if err != nil {
		return contract.TaskDashboardKPI{}, contract.TopExecutorByCompleted{},
			fmt.Errorf("get top executor by completed in period: %w", err)
	}

	return kpi, top, nil
}

// normalizeAnalyticsPeriod: если границы не заданы — последние 30 суток до текущего UTC-момента
//
// если задана только to — from = to − 30 дней
//
// если только from — to = сейчас.
func normalizeAnalyticsPeriod(from, to time.Time) (time.Time, time.Time, error) {
	now := time.Now().UTC()

	switch {
	case to.IsZero() && from.IsZero():
		to = now
		from = to.AddDate(0, 0, -30)
	case to.IsZero():
		to = now
	case from.IsZero():
		from = to.AddDate(0, 0, -30)
	}

	if from.After(to) {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: from after to", contract.ErrInvalidAnalyticsPeriod)
	}

	return from.UTC(), to.UTC(), nil
}
