package memo

import (
	"context"
	"memoapi/model"
)

// Repository defines memo repository contract.
type Repository interface {
	CreateMemo(ctx context.Context, a *model.Memo) error
	GetMemoByID(ctx context.Context, id int) (*model.Memo, error)
	GetAllMemo(ctx context.Context) ([]*model.Memo, error)
}
