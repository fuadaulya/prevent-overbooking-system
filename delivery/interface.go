package delivery

import "context"

type Usecase interface {
	AddItemToCart(ctx context.Context, cartID, userID, productID, quantity int) error
	Checkout(ctx context.Context, cartID, userID int) error
	ValidateCartOwnership(ctx context.Context, cartID, userID int) bool
}
