package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">>before:", account1.Balance, account2.Balance)

	//concurrency transfer transaction testing
	n := 5
	amount := int64(2)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {

		txName := fmt.Sprintf("tx %d", i+1)

		go func() {

			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			if err != nil {
				errs <- err
			} else {
				results <- result
			}
		L:
			for {
				select {
				case <-errs:
					fmt.Println("error:", err)
					break L

				case <-results:
					fmt.Println("results:", result)
					break L

				default:
					break
				}
			}

		}()

		time.Sleep(1 * time.Second)

	}
	//checking results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		//require.NoError(t, err)

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

		//checking accounts
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		//checking balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := account2.Balance - account2.Balance

		require.Equal(t, diff1, diff2) //they should be equal
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount, 2*amount, 3*amount, ..... n*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= int(n))

		require.NotContains(t, existed, k)
		existed[k] = true //avoid double transaction bp
	}

	//check the final balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Println("stopped")
	close(errs)
	<-errs
	close(results)
	<-results

	return

}
