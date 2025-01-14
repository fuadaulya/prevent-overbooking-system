package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"task-one/pkg/entity"
)

// GetAllProducts retrieve all product data
func (r *Postgres) GetAllProducts(ctx context.Context) ([]entity.Products, error) {
	var urls []entity.Products

	query := `SELECT
				id,
				name,
				price,
				stock
			FROM
				products`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to get all urls: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var url entity.Products
		err := rows.Scan(
			&url.ID,
			&url.Name,
			&url.Price,
			&url.Stock,
		)
		if err != nil {
			log.Printf("Failed to scan url: %+v", err)
			return nil, fmt.Errorf("[postgres] failed to get all urls: %+v", err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (r *Postgres) GetProductByID(ctx context.Context, productID int) (entity.Products, error) {
	var product entity.Products

	query := `SELECT
				id,
				name,
				stock,
				price
			FROM
				products
			WHERE id = $1`

	err := r.DB.QueryRowContext(ctx, query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Stock,
		&product.Price,
	)

	if err != nil {
		return product, fmt.Errorf("error retrieving product: %+v", err)
	}
	return product, nil
}

// GetStock retrieve the current stock of a product
func (r *Postgres) GetStock(ctx context.Context, productID int) (int, error) {
	var stock int
	err := r.DB.QueryRowContext(ctx, "SELECT stock FROM products WHERE id = $1", productID).Scan(&stock)
	if err != nil {
		return 0, fmt.Errorf("[postgres] error when retrieving the current stock: %+v", err)
	}
	return stock, nil
}

// ReduceStock reduces the stock for a product if stock is available
func (r *Postgres) ReduceStock(ctx context.Context, productID int, quantity int) error {
	// Start transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[postgres] failed to start transaction: %+v", err)
	}
	defer tx.Rollback() // Rollback if failed

	// Retrieve the current stock with a transaction lock
	var stock int
	err = tx.QueryRowContext(ctx, "SELECT stock FROM products WHERE id = $1 FOR UPDATE", productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("[postgres] failed to get stock: %+v", err)
	}

	// Validate if the stock is sufficient
	if stock < quantity {
		return errors.New("not enough stock")
	}

	// Reduce the product stock
	_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", quantity, productID)
	if err != nil {
		return fmt.Errorf("[postgres] failed to execute query: %+v", err)
	}

	// Commit the transaction if all operations succeed
	return tx.Commit()
}

func (r *Postgres) ReserveStock(ctx context.Context, productID, quantity int) error {
	// Start a transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[postgres] failed to start transaction: %+v", err)
	}
	defer tx.Rollback() // Ensure rollback on failure

	// Check product stock with a lock
	var stock int
	queryGetStock := `SELECT
						stock
					FROM
						products
					WHERE
						id = $1 FOR UPDATE`

	err = tx.QueryRowContext(ctx, queryGetStock, productID).Scan(&stock)
	if err != nil {
		return fmt.Errorf("[postgres] failed to get stock: %+v", err)
	}

	// Validate if stock is sufficient
	if stock < quantity {
		return errors.New("not enough stock")
	}

	// Reduce stock in the products table
	queryReduceStock := `UPDATE
							products
						SET
							stock = stock - $1
						WHERE
							id = $2`

	_, err = tx.ExecContext(ctx, queryReduceStock, quantity, productID)
	if err != nil {
		return fmt.Errorf("[postgres] failed to reduce stock: %+v", err)
	}

	// Insert reservation into reserved_stock
	queryReservation := `INSERT INTO
							reserved_stock (product_id, quantity, created_at, updated_at, status)
						VALUES
							($1, $2, NOW(), NOW(), 'reserved')`

	_, err = tx.ExecContext(ctx, queryReservation, productID, quantity)
	if err != nil {
		return fmt.Errorf("[postgres] failed to insert reserved stock: %+v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("[postgres] failed to commit transaction: %+v", err)
	}

	return nil
}
