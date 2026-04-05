package office

import (
	"context"
	"fmt"

	"github.com/gigasigmaslav/guard-panel-api/internal/domain/entity"
)

type soOfficeRepo interface {
	GetOffices(ctx context.Context) ([]entity.Office, error)
}

type SearchOfficeUseCase struct {
	officeRepo soOfficeRepo
}

func NewSearchOfficeUseCase(
	officeRepo soOfficeRepo,
) *SearchOfficeUseCase {
	return &SearchOfficeUseCase{
		officeRepo: officeRepo,
	}
}

func (so *SearchOfficeUseCase) Search(ctx context.Context) ([]entity.Office, error) {
	offices, err := so.officeRepo.GetOffices(ctx)
	if err != nil {
		return nil, fmt.Errorf("search offices: %w", err)
	}
	return offices, nil
}
