package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	auth "github.com/gigasigmaslav/guard-panel-api/internal/pkg/auth"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/contract"
	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type siUserRepo interface {
	GetUserByEmployeeID(ctx context.Context, employeeID int64) (entity.User, error)
}

type siEmployeeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
}

type SignInUseCase struct {
	userRepo     siUserRepo
	employeeRepo siEmployeeRepo
	codec        auth.TokenCodec
	hasher       auth.PasswordHasher
}

func NewSignInUseCase(
	userRepo siUserRepo,
	employeeRepo siEmployeeRepo,
	codec auth.TokenCodec,
	hasher auth.PasswordHasher,
) *SignInUseCase {
	return &SignInUseCase{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		codec:        codec,
		hasher:       hasher,
	}
}

func (si *SignInUseCase) SignIn(
	ctx context.Context,
	employeeID int64,
	plainPassword string,
) (contract.AuthResult, error) {
	u, getUserErr := si.userRepo.GetUserByEmployeeID(ctx, employeeID)
	if getUserErr != nil {
		if errors.Is(getUserErr, entity.ErrUserNotFound) {
			return contract.AuthResult{}, contract.ErrInvalidCredentials
		}

		return contract.AuthResult{}, fmt.Errorf("sign in get user: %w", getUserErr)
	}

	if compareErr := si.hasher.Compare(u.PasswordHash, plainPassword); compareErr != nil {
		return contract.AuthResult{}, contract.ErrInvalidCredentials
	}

	emp, getEmployeeErr := si.employeeRepo.GetEmployeeByID(ctx, employeeID)
	if getEmployeeErr != nil {
		if errors.Is(getEmployeeErr, entity.ErrEmployeeNotFound) {
			return contract.AuthResult{}, contract.ErrInvalidCredentials
		}

		return contract.AuthResult{}, fmt.Errorf("sign in get employee: %w", getEmployeeErr)
	}

	token, exp, issueErr := si.codec.IssueBearer(employeeID, int32(emp.Position), time.Now())
	if issueErr != nil {
		return contract.AuthResult{}, fmt.Errorf("sign in issue token: %w", issueErr)
	}

	return contract.AuthResult{
		AccessToken: token,
		ExpiresAt:   exp,
	}, nil
}
