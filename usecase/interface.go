package usecase

import (
	"context"
	"task-one/pkg/entity"

	db "task-one/repository/database"
)

type DBRepo interface {
	// Cart
	CheckCartExists(ctx context.Context, cartID int) (bool, error)
	CreateCart(ctx context.Context, userID int) (int, error)
	GetCartItemsByCartID(ctx context.Context, cartID int) ([]entity.CartItem, error)
	AddItem(ctx context.Context, cartID, productID, quantity int) error
	GetCartOwner(ctx context.Context, cartID int) (int, error)

	// Transaction
	BeginTransaction(ctx context.Context) (*db.DBTransaction, error)
}

type RedisRepo interface {
	GetProductStockFromCache(ctx context.Context, productID int) (int, error)
	SetProductStockInCache(ctx context.Context, productID, stock int) error
}
