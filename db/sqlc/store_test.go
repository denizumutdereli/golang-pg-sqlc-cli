package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//concurrency transfer transaction testing
	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			if err != nil {
				errors <- err
			}

			results <- result

		}()
	}

	//checking results
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//checking transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check orders
		fromOrder := result.FromOrder
		require.NotEmpty(t, fromOrder)
		require.Equal(t, account1.ID, fromOrder.AccountID)
		require.Equal(t, -amount, fromOrder.Amount) //deduction
		require.NotZero(t, fromOrder.ID)
		require.NotZero(t, fromOrder.CreatedAt)

		_, err = store.GetOrder(context.Background(), fromOrder.ID)
		require.NoError(t, err)

		toOrder := result.ToOrder
		require.NotEmpty(t, toOrder)
		require.Equal(t, account1.ID, toOrder.AccountID)
		require.Equal(t, amount, toOrder.Amount) //deduction
		require.NotZero(t, toOrder.ID)
		require.NotZero(t, toOrder.CreatedAt)

		_, err = store.GetOrder(context.Background(), toOrder.ID)
		require.NoError(t, err)

		//checking result balances. I will back here...

	}

}
