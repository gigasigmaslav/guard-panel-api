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

type suEmployeeRepo interface {
	GetEmployeeByID(ctx context.Context, id int64) (entity.Employee, error)
}

type SignUpUseCase struct {
	employeeRepo suEmployeeRepo
	transactor   contract.RepoTransactor
	codec        auth.TokenCodec
	hasher       auth.PasswordHasher
}

func NewSignUpUseCase(
	employeeRepo suEmployeeRepo,
	transactor contract.RepoTransactor,
	codec auth.TokenCodec,
	hasher auth.PasswordHasher,
) *SignUpUseCase {
	return &SignUpUseCase{
		employeeRepo: employeeRepo,
		transactor:   transactor,
		codec:        codec,
		hasher:       hasher,
	}
}

func (su *SignUpUseCase) SignUp(
	ctx context.Context,
	employeeID int64,
	plainPassword string,
) (contract.AuthResult, error) {
	emp, getEmployeeErr := su.employeeRepo.GetEmployeeByID(ctx, employeeID)
	if getEmployeeErr != nil {
		if errors.Is(getEmployeeErr, entity.ErrEmployeeNotFound) {
			return contract.AuthResult{}, fmt.Errorf("sign up: %w", getEmployeeErr)
		}
		return contract.AuthResult{}, fmt.Errorf("sign up get employee: %w", getEmployeeErr)
	}

	if emp.Position == entity.EmployeePositionUnspecified {
		return contract.AuthResult{}, contract.ErrInvalidEmployeePosition
	}

	hash, hashErr := su.hasher.Hash(plainPassword)
	if hashErr != nil {
		return contract.AuthResult{}, fmt.Errorf("sign up hash password: %w", hashErr)
	}

	if txErr := su.createUserInTx(ctx, employeeID, hash); txErr != nil {
		if errors.Is(txErr, contract.ErrAlreadyRegistered) {
			return contract.AuthResult{}, contract.ErrAlreadyRegistered
		}
		return contract.AuthResult{}, fmt.Errorf("sign up: %w", txErr)
	}

	token, exp, issueErr := su.codec.IssueBearer(employeeID, int32(emp.Position), time.Now())
	if issueErr != nil {
		return contract.AuthResult{}, fmt.Errorf("sign up issue token: %w", issueErr)
	}

	return contract.AuthResult{
		AccessToken: token,
		ExpiresAt:   exp,
	}, nil
}

func (su *SignUpUseCase) createUserInTx(
	ctx context.Context,
	employeeID int64,
	passwordHash string,
) error {
	return su.transactor.InTx(ctx, func(tx contract.TxRepo) error {
		userExists, existErr := tx.UserExistsByEmployeeID(ctx, employeeID)
		if existErr != nil {
			return fmt.Errorf("sign up check user exists: %w", existErr)
		}
		if userExists {
			return contract.ErrAlreadyRegistered
		}
		if _, createErr := tx.CreateUser(ctx, entity.User{
			EmployeeID:   employeeID,
			PasswordHash: passwordHash,
		}); createErr != nil {
			return fmt.Errorf("sign up create user: %w", createErr)
		}
		return nil
	})
}
