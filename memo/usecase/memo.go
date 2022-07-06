package usecase

import (
	"context"
	"memoapi/memo"
	"memoapi/model"

	"github.com/go-playground/validator/v10"
)

// MemoUsecase implements memo usecase.
type MemoUsecase struct {
	repository memo.Repository
	validate   *validator.Validate
}

// NewMemoUsecase returns new instance of Usecase.
func NewMemoUsecase(repository memo.Repository) *MemoUsecase {
	return &MemoUsecase{
		repository: repository,
		validate:   validator.New(),
	}
}

// CreateMemo creates a new memo.
func (u *MemoUsecase) CreateMemo(ctx context.Context, m *model.Memo) error {
	if err := u.validate.Struct(m); err != nil {
		return err
	}
	return u.repository.CreateMemo(ctx, m)
}

// GetMemoByID gets memo by given id.
func (u *MemoUsecase) GetMemoByID(ctx context.Context, id int) (*model.Memo, error) {
	return u.repository.GetMemoByID(ctx, id)
}

// GetAllMemo gets all memos stored in repository.
func (u *MemoUsecase) GetAllMemo(ctx context.Context) ([]*model.Memo, error) {
	return u.repository.GetAllMemo(ctx)
}
