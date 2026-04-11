package server

import (
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/auth"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/comment"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/employee"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/office"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/refund"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/task"
	vuddecision "github.com/gigasigmaslav/guard-panel-api/internal/domain/usecase/vud-decision"
)

type Dependencies struct {
	*comment.CreateCommentUseCase
	*comment.DeleteCommentUseCase
	*refund.CreateRefundUseCase
	*vuddecision.CreateVUDDecisionUseCase
	*vuddecision.UpdateVUDDecisionUseCase
	*employee.CreateEmployeeUseCase
	*employee.UpdateEmployeeUseCase
	*employee.DeleteEmployeeUseCase
	*employee.SearchEmployeeUseCase
	*office.CreateOfficeUseCase
	*office.UpdateOfficeUseCase
	*office.DeleteOfficeUseCase
	*office.SearchOfficeUseCase
	*task.CreateTaskUseCase
	*task.UpdateTaskUseCase
	*task.GetTaskByIDUseCase
	*task.SearchTasksUseCase
	*auth.SignUpUseCase
	*auth.SignInUseCase
	*auth.WhoAmIUseCase
}

func NewDependencies(
	createCommentUC *comment.CreateCommentUseCase,
	deleteCommentUC *comment.DeleteCommentUseCase,
	createRefundUC *refund.CreateRefundUseCase,
	createVUDDecisionUC *vuddecision.CreateVUDDecisionUseCase,
	updateVUDDecisionUC *vuddecision.UpdateVUDDecisionUseCase,
	createEmployeeUC *employee.CreateEmployeeUseCase,
	updateEmployeeUC *employee.UpdateEmployeeUseCase,
	deleteEmployeeUC *employee.DeleteEmployeeUseCase,
	searchEmployeeUC *employee.SearchEmployeeUseCase,
	createOfficeUC *office.CreateOfficeUseCase,
	updateOfficeUC *office.UpdateOfficeUseCase,
	deleteOfficeUC *office.DeleteOfficeUseCase,
	searchOfficeUC *office.SearchOfficeUseCase,
	createTaskUC *task.CreateTaskUseCase,
	updateTaskUC *task.UpdateTaskUseCase,
	getTaskByIDUC *task.GetTaskByIDUseCase,
	searchTasksUC *task.SearchTasksUseCase,
	signUpUC *auth.SignUpUseCase,
	signInUC *auth.SignInUseCase,
	whoAmIUC *auth.WhoAmIUseCase,
) *Dependencies {
	return &Dependencies{
		CreateCommentUseCase:     createCommentUC,
		DeleteCommentUseCase:     deleteCommentUC,
		CreateRefundUseCase:      createRefundUC,
		CreateVUDDecisionUseCase: createVUDDecisionUC,
		UpdateVUDDecisionUseCase: updateVUDDecisionUC,
		CreateEmployeeUseCase:    createEmployeeUC,
		UpdateEmployeeUseCase:    updateEmployeeUC,
		DeleteEmployeeUseCase:    deleteEmployeeUC,
		SearchEmployeeUseCase:    searchEmployeeUC,
		CreateOfficeUseCase:      createOfficeUC,
		UpdateOfficeUseCase:      updateOfficeUC,
		DeleteOfficeUseCase:      deleteOfficeUC,
		SearchOfficeUseCase:      searchOfficeUC,
		CreateTaskUseCase:        createTaskUC,
		UpdateTaskUseCase:        updateTaskUC,
		GetTaskByIDUseCase:       getTaskByIDUC,
		SearchTasksUseCase:       searchTasksUC,
		SignUpUseCase:            signUpUC,
		SignInUseCase:            signInUC,
		WhoAmIUseCase:            whoAmIUC,
	}
}
