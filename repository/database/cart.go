package repository

import (
	"context"
	"database/sql"
	"fmt"
	"task-one/pkg/entity"
)

// Check if the cart already exists
func (r *Postgres) CheckCartExists(ctx context.Context, cartID int) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM carts WHERE id = $1)", cartID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("[postgres] failed to check cart existence: %+v", err)
	}
	return exists, nil
}

// Create a new cart if it doesn't exist
func (r *Postgres) CreateCart(ctx context.Context, userID int) (int, error) {
	var cartID int

	query := `INSERT INTO carts
				(user_id, created_at, updated_at)
        	VALUES
				($1, NOW(), NOW())
        	RETURNING id`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&cartID)
	if err != nil {
		return 0, fmt.Errorf("[postgres] failed to create cart: %+v", err)
	}

	return cartID, nil
}

func (r *Postgres) AddItem(ctx context.Context, cartID, productID, quantity int) error {
	// Check if the product already exists in the cart
	existingQuantity, err := r.getProductQuantityInCart(ctx, cartID, productID)
	if err != nil {
		return err
	}

	// If the product exists, update the quantity, otherwise insert a new item
	if existingQuantity > 0 {
		return r.updateCartItemQuantity(ctx, cartID, productID, quantity)
	}
	return r.insertCartItem(ctx, cartID, productID, quantity)
}

// Get the existing quantity of a product in the cart
func (r *Postgres) getProductQuantityInCart(ctx context.Context, cartID, productID int) (int, error) {
	var existingQuantity int

	query := `UPDATE
				reserved_stock
			SET
				status = $1
			WHERE
				product_id = $2`

	err := r.DB.QueryRowContext(ctx, query, cartID, productID).Scan(&existingQuantity)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("[postgres] failed to check if product exists in cart: %+v", err)
	}

	return existingQuantity, nil
}

// Update the quantity of an existing item in the cart
func (r *Postgres) updateCartItemQuantity(ctx context.Context, cartID, productID, quantity int) error {
	query := `UPDATE
				cart_items
			SET
				quantity = quantity + $3, updated_at = NOW()
			WHERE
				cart_id = $1 AND product_id = $2`

	_, err := r.DB.ExecContext(ctx, query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("[postgres] failed to update cart item: %+v", err)
	}

	return nil
}

// Insert a new item into the cart
func (r *Postgres) insertCartItem(ctx context.Context, cartID, productID, quantity int) error {
	query := `INSERT INTO
				cart_items (cart_id, product_id, quantity, created_at, updated_at)
			VALUES
				($1, $2, $3, NOW(), NOW())`

	_, err := r.DB.ExecContext(ctx, query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("[postgres] failed to add item to cart: %+v", err)
	}

	return nil
}

// Function to retrieve ordered products in the cart based on cart_id
func (r *Postgres) GetCartItemsByCartID(ctx context.Context, cartID int) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem

	query := `SELECT DISTINCT
				ci.product_id, ci.quantity, rs.status
			FROM
				cart_items ci
			LEFT JOIN
				reserved_stock rs ON ci.product_id = rs.product_id
			WHERE
				ci.cart_id = $1
			AND
				rs.status = 'reserved'`

	rows, err := r.DB.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to query cart_items: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.CartItem
		err := rows.Scan(&item.ProductID,
			&item.Quantity,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("[postgres] failed to scan cart_items row: %+v", err)
		}

		cartItems = append(cartItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[postgres] error while iterating rows: %+v", err)
	}

	return cartItems, nil
}

// Repository function to get the owner of a cart based on cart_id
func (r *Postgres) GetCartOwner(ctx context.Context, cartID int) (int, error) {
	var ownerID int

	// Query to get the user_id who owns the cart
	query := `SELECT 
				user_id
			FROM 
				carts
			WHERE
				id = $1`

	err := r.DB.QueryRowContext(ctx, query, cartID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("cart not found")
		}
		return 0, fmt.Errorf("failed to query cart owner: %v", err)
	}

	return ownerID, nil
}
