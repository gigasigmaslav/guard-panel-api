package contract

import (
	"time"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
	"github.com/gigasigmaslav/guard-panel-api/pkg/api/v1/message"
)

type SearchTasksRequest struct {
	Status  entity.TaskStatus
	Page    int32
	PerPage int32
	Filter  *SearchTasksFilter
	Sorting SearchTasksSorting
}

type SearchTasksFilter struct {
	ID           *int64
	Priority     *entity.TaskPriority
	OfficeID     *int64
	ExecutorID   *int64
	ViolatorType *entity.ViolatorType
	KUSP         *string
	UD           *string
}

type SearchTasksSorting struct {
	OrderDirection OrderDirection
	OrderBy        OrderBy
}

type OrderDirection int32

const (
	OrderDirectionUnspecified OrderDirection = iota
	OrderDirectionAsc
	OrderDirectionDesc
)

type OrderBy int32

const (
	OrderByUnspecified OrderBy = iota
	OrderByDamageAmountKopecks
	OrderByCreatedAt
)

type TaskLookupItem struct {
	ID            int64
	Status        entity.TaskStatus
	DamageAmount  int64
	Priority      entity.TaskPriority
	ExecutorID    int64
	ExecutorName  string
	ViolatorID    int64
	ViolatorName  string
	OfficeID      int64
	OfficeName    string
	CreatedByID   int64
	CreatedByName string
	CreatedAt     time.Time
	StartDate     time.Time
	EndDate       *time.Time
}

type TaskDetailsDTO struct {
	Task          TaskLookupItem
	Refunds       []entity.Refund
	VUDDecision   *entity.VUDDecision
	Comments      []entity.Comment
	HistoryEvents []entity.TaskHistoryChange
}

func MapToSearchTasksFilter(filter *message.SearchTasksRequest_Filter) *SearchTasksFilter {
	if filter == nil {
		return nil
	}

	result := &SearchTasksFilter{}

	if priority := mapTaskPriorityToContractPtr(filter.GetPriority()); priority != nil {
		result.Priority = priority
	}

	if officeID := filter.GetOfficeId(); officeID > 0 {
		result.OfficeID = &officeID
	}

	if executorID := filter.GetExecutorId(); executorID > 0 {
		result.ExecutorID = &executorID
	}

	if violatorType := mapViolatorTypeToContractPtr(filter.GetViolatorType()); violatorType != nil {
		result.ViolatorType = violatorType
	}

	if kusp := filter.GetKusp(); kusp != "" {
		result.KUSP = &kusp
	}

	if ud := filter.GetUd(); ud != "" {
		result.UD = &ud
	}

	return result
}

func MapTaskStatusToContract(status message.TaskStatus) entity.TaskStatus {
	switch status {
	case message.TaskStatus_CASE_STATUS_UNSPECIFIED:
		return entity.TaskStatusUnspecified
	case message.TaskStatus_NEW:
		return entity.TaskStatusNew
	case message.TaskStatus_DOCUMENT_GATHERING:
		return entity.TaskStatusDocumentGathering
	case message.TaskStatus_WAITING_FOR_VUD:
		return entity.TaskStatusWaitingForVUD
	case message.TaskStatus_JUDICIAL_PROCEEDINGS:
		return entity.TaskStatusJudicialProceedings
	case message.TaskStatus_COMPLETED:
		return entity.TaskStatusCompleted
	default:
		return entity.TaskStatusUnspecified
	}
}

func MapOrderDirectionToContract(direction message.OrderDirection) OrderDirection {
	switch direction {
	case message.OrderDirection_ORDER_DIRECTION_UNSPECIFIED:
		return OrderDirectionUnspecified
	case message.OrderDirection_ASC:
		return OrderDirectionAsc
	case message.OrderDirection_DESC:
		return OrderDirectionDesc
	default:
		return OrderDirectionUnspecified
	}
}

func MapOrderByToContract(orderBy message.OrderBy) OrderBy {
	switch orderBy {
	case message.OrderBy_ORDER_BY_UNSPECIFIED:
		return OrderByUnspecified
	case message.OrderBy_BARCODES_DAMAGE_AMOUNT:
		return OrderByDamageAmountKopecks
	case message.OrderBy_CREATED_AT:
		return OrderByCreatedAt
	default:
		return OrderByUnspecified
	}
}

func mapTaskPriorityToContractPtr(priority message.TaskPriority) *entity.TaskPriority {
	var p entity.TaskPriority

	switch priority {
	case message.TaskPriority_TASK_PRIORITY_UNSPECIFIED:
		return nil
	case message.TaskPriority_LOW:
		p = entity.TaskPriorityLow
	case message.TaskPriority_HIGH:
		p = entity.TaskPriorityHigh
	default:
		return nil
	}

	return &p
}

func mapViolatorTypeToContractPtr(t message.ViolatorType) *entity.ViolatorType {
	var violatorType entity.ViolatorType

	switch t {
	case message.ViolatorType_VIOLATOR_TYPE_UNSPECIFIED:
		return nil
	case message.ViolatorType_EMPLOYEE:
		violatorType = entity.ViolatorTypeEmployee
	case message.ViolatorType_CLIENT:
		violatorType = entity.ViolatorTypeClient
	default:
		return nil
	}

	return &violatorType
}
