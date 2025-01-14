package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func (u *Usecase) EnsureCartExists(ctx context.Context, cartID, userID int) (int, error) {
	// Check if the cart already exists
	exists, err := u.dbrepo.CheckCartExists(ctx, cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to check if cart exists: %+v", err)
	}

	if !exists {
		// If not, create a new cart
		newCartID, err := u.dbrepo.CreateCart(ctx, userID)
		if err != nil {
			return 0, fmt.Errorf("failed to create new cart: %+v", err)
		}
		return newCartID, nil
	}

	// If the cart exists, return the existing cartID
	return cartID, nil
}

func (u *Usecase) ReserveStock(ctx context.Context, productID, quantity int) error {
	// Get stock from Redis
	stock, err := u.redisrepo.GetProductStockFromCache(ctx, productID)
	if err != nil && err != redis.Nil {
		return err
	}

	// If stock is not in cache, retrieve it from the database
	// Start a transaction
	tx, err := u.dbrepo.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Ensure rollback on failure

	// Get stock from database
	if stock == 0 {
		stock, err = tx.GetStockWithLock(ctx, productID)
		if err != nil {
			return err
		}
	}

	// Validate product stock
	if quantity <= 0 || quantity > stock {
		return fmt.Errorf("invalid quantity or insufficient stock")
	}

	// Reduce stock in database and get the new stock
	newStock, err := tx.ReduceStock(ctx, productID, quantity)
	if err != nil {
		return err
	}

	// Insert reservation into reserved_stock
	err = tx.InsertReservation(ctx, productID, quantity)
	if err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %+v", err)
	}

	// Update stock to Redis
	err = u.redisrepo.SetProductStockInCache(ctx, productID, newStock)
	if err != nil {
		log.Print(err)
	}

	return nil
}

func (u *Usecase) AddItemToCart(ctx context.Context, cartID, userID, productID, quantity int) error {
	// Ensure the cart exists
	cartID, err := u.EnsureCartExists(ctx, cartID, userID)
	if err != nil {
		return err
	}

	// Check and reduce stock in product repository
	err = u.ReserveStock(ctx, productID, quantity)
	if err != nil {
		return err
	}

	// Add item to cart
	err = u.dbrepo.AddItem(ctx, cartID, productID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Checkout(ctx context.Context, cartID, userID int) error {
	// Start transaction
	tx, err := u.dbrepo.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %+v", err)
	}
	defer tx.Rollback()

	// Validate and confirm reserved stock
	items, err := u.dbrepo.GetCartItemsByCartID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart items: %+v", err)
	}

	for _, item := range items {
		if item.Status != "reserved" {
			return fmt.Errorf("reserved stock not valid for product %d", item.ProductID)
		}

		// Update reserved stock
		err = tx.UpdateReservedStockStatus(ctx, item.ProductID, "confirmed")
		if err != nil {
			return fmt.Errorf("failed to confirm reserved stock for product %d: %+v", item.ProductID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %+v", err)
	}

	return nil
}

/*
	HandleRestoreStock function can be used in the future for several scenarios, including:

1. After receiving a failed response from the payment gateway
2. Timeout from the internal system
  - If the user starts the payment process but does not complete it within a certain time
    (e.g., the user closes the app or does not continue the payment after being redirected to the payment page).

3. Batch System or Scheduler
  - To handle timeout cases or orders that have been pending for too long,
    the application can use a batch job or cron job to check for orders with a "reserved" status that have not been paid within a certain period.
*/
func (u *Usecase) HandleRestoreStock(ctx context.Context, cartID int) error {
	// Start transaction
	tx, err := u.dbrepo.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Retrieve reserved stock for the cart
	reservedItems, err := u.dbrepo.GetCartItemsByCartID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get reserved items for cart %d: %v", cartID, err)
	}

	for _, item := range reservedItems {
		// Update reserved_stock status to 'cancelled'
		err := tx.UpdateReservedStockStatus(ctx, item.ProductID, "cancelled")
		if err != nil {
			return fmt.Errorf("failed to update reserved stock status for product %d: %v", item.ProductID, err)
		}

		// Return reserved stock back to products
		err = tx.RestoreStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to restore stock for product %d: %v", item.ProductID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Function to verify if the cart belongs to the provided user.
func (u *Usecase) ValidateCartOwnership(ctx context.Context, cartID, userID int) bool {
	// Call the repository function to check if the cart belongs to the user
	ownerID, err := u.dbrepo.GetCartOwner(ctx, cartID)
	if err != nil {
		// If an error occurs, log or handle the error
		return false
	}

	// Compare the ownerID with the provided userID
	return ownerID == userID
}
