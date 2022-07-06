package memo

import (
	"context"
	"memoapi/model"
)

// Usecase defines memo usecase contract.
type Usecase interface {
	CreateMemo(ctx context.Context, m *model.Memo) error
	GetMemoByID(ctx context.Context, id int) (*model.Memo, error)
	GetAllMemo(ctx context.Context) ([]*model.Memo, error)
}
