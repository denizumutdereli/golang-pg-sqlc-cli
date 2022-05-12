package db

import (
	"context"
	"database/sql"
	"fmt"
)

/*
Composition DB Transaction
*/

// store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rollBackError := tx.Rollback(); rollBackError != nil {
			return fmt.Errorf("tx err:%w", rollBackError)
		}
		return err
	}

	//all done no problem
	return tx.Commit()
}

/***
Checking...
Lets transfer money from one account to an another one
This going to create a transfer instance and account orders and update the accounts balance with in the same db transaction block.
***/

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromOrder   Order    `json:"from_entry`
	ToOrder     Order    `json:"to_entry`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromOrder, err = q.CreateOrder(ctx, CreateOrderParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		result.ToOrder, err = q.CreateOrder(ctx, CreateOrderParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//... update the balance.. I will back to here later..

		return nil
	})

	return result, err

}
