package comment

import (
	"context"
	"fmt"
)

type dcCommentRepo interface {
	DeleteCommentByID(ctx context.Context, id int64) error
}

type DeleteCommentUseCase struct {
	commentRepo dcCommentRepo
}

func NewDeleteCommentUseCase(
	commentRepo dcCommentRepo,
) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{
		commentRepo: commentRepo,
	}
}

func (dc *DeleteCommentUseCase) Delete(ctx context.Context, id int64) error {
	if err := dc.commentRepo.DeleteCommentByID(ctx, id); err != nil {
		return fmt.Errorf("delete comment by id: %w", err)
	}
	return nil
}
