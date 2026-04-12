package contract

import "errors"

var ErrInvalidAnalyticsPeriod = errors.New("invalid analytics period")

// TaskDashboardKPI агрегаты для дашборда: активные задачи «сейчас» и метрики за [from, to].
type TaskDashboardKPI struct {
	ActiveTasksCount        int64
	CompletedInPeriod       int64
	CreatedInPeriod         int64
	CompletedToCreatedRatio *float64
}

// TopExecutorByCompleted топ исполнитель по числу задач, закрытых за период (по end_date).
type TopExecutorByCompleted struct {
	ExecutorID   int64
	ExecutorName string
	Completed    int64
}
