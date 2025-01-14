package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTransaction struct {
	tx  *sql.Tx
	db  *Postgres
	ctx context.Context
}

// Begin a new transaction
func (r *Postgres) BeginTransaction(ctx context.Context) (*DBTransaction, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to begin transaction: %+v", err)
	}
	return &DBTransaction{
		tx:  tx,
		db:  r,
		ctx: ctx,
	}, nil
}

// Commit the transaction
func (t *DBTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback the transaction
func (t *DBTransaction) Rollback() error {
	return t.tx.Rollback()
}

// Check product stock with a lock
func (t *DBTransaction) GetStockWithLock(ctx context.Context, productID int) (int, error) {
	var stock int

	query := `SELECT
				stock
			FROM
				products
			WHERE
				id = $1 FOR UPDATE`

	err := t.tx.QueryRowContext(t.ctx, query, productID).Scan(&stock)
	if err != nil {
		return stock, fmt.Errorf("[postgres] failed to get stock: %+v", err)
	}

	return stock, nil
}

// Reduce stock in the products and return the new stock
func (t *DBTransaction) ReduceStock(ctx context.Context, productID, quantity int) (int, error) {
	var newStock int

	query := `UPDATE
				products
			SET
				stock = stock - $1
			WHERE
				id = $2
			RETURNING stock`

	err := t.tx.QueryRowContext(ctx, query, quantity, productID).Scan(&newStock)
	if err != nil {
		return 0, fmt.Errorf("[postgres] failed to reduce stock: %+v", err)
	}

	return newStock, nil
}

// Insert reservation into reserved_stock
func (t *DBTransaction) InsertReservation(ctx context.Context, productID, quantity int) error {
	query := `INSERT INTO
				reserved_stock (product_id, quantity, created_at, updated_at, status)
			VALUES
				($1, $2, NOW(), NOW(), 'reserved')`

	_, err := t.tx.ExecContext(t.ctx, query, productID, quantity)
	if err != nil {
		return fmt.Errorf("[postgres] failed to insert reserved stock: %+v", err)
	}
	return nil
}

// Update reserved stock status
func (t *DBTransaction) UpdateReservedStockStatus(ctx context.Context, productID int, status string) error {
	query := `UPDATE
				reserved_stock
			SET
				status = $1
			WHERE
				product_id = $2`

	_, err := t.tx.ExecContext(ctx, query, status, productID)
	if err != nil {
		return fmt.Errorf("[postgres] failed to update status: %+v", err)
	}

	return nil
}

// Function to return the ordered product stock
func (t *DBTransaction) RestoreStock(ctx context.Context, productID, quantity int) error {
	// Restoring stock in the products table
	updateStockQuery := `UPDATE
							products
						SET
							stock = stock + $1
						WHERE
							id = $2`

	_, err := t.tx.ExecContext(ctx, updateStockQuery, quantity, productID)
	if err != nil {
		return fmt.Errorf("[postgres] failed to update product stock: %+v", err)
	}

	return nil
}
