package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/chutified/simple-bank/util"
)

func TestStore_TransferTx(t *testing.T) {
	s := db.NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// test concurrent transfer transactions
	n := 10
	arg := db.TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(),
	}

	results := make(chan db.TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := s.TransferTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// transfer
		transfer := result.Transfer
		if assert.NotEmpty(t, transfer) {
			assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
			assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
			assert.Equal(t, arg.Amount, transfer.Amount)

			assert.NotZero(t, transfer.ID)
			assert.NotZero(t, transfer.CreatedAt)
		}

		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// entryFrom
		entryFrom := result.FromEntry
		if assert.NotEmpty(t, entryFrom) {
			assert.Equal(t, arg.FromAccountID, entryFrom.AccountID)
			assert.Equal(t, arg.Amount, -entryFrom.Amount)

			assert.NotZero(t, entryFrom.ID)
			assert.NotZero(t, entryFrom.CreatedAt)
		}

		_, err = testQueries.GetEntry(context.Background(), transfer.FromAccountID)
		require.NoError(t, err)

		// entryTo
		entryTo := result.ToEntry
		if assert.NotEmpty(t, entryTo) {
			assert.Equal(t, arg.ToAccountID, entryTo.AccountID)
			assert.Equal(t, arg.Amount, entryTo.Amount)

			assert.NotZero(t, entryTo.ID)
			assert.NotZero(t, entryTo.CreatedAt)
		}

		_, err = testQueries.GetEntry(context.Background(), transfer.ToAccountID)
		require.NoError(t, err)

		// accounts
		fromAccount := result.FromAccount
		if assert.NotEmpty(t, fromAccount) {
			assert.Equal(t, arg.FromAccountID, fromAccount.ID)
		}

		_, err = testQueries.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)

		toAccount := result.ToAccount
		if assert.NotEmpty(t, toAccount) {
			assert.Equal(t, arg.ToAccountID, toAccount.ID)
		}

		_, err = testQueries.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)

		// accounts' balances
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := account2.Balance - toAccount.Balance

		assert.Equal(t, arg.Amount*int64(i+1), diff1)
		assert.Equal(t, -arg.Amount*int64(i+1), diff2)
	}

	// check final accounts
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), arg.FromAccountID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), arg.ToAccountID)
	require.NoError(t, err)

	assert.Equal(t, account1.Balance-arg.Amount*int64(n), updatedFromAccount.Balance)
	assert.Equal(t, account2.Balance+arg.Amount*int64(n), updatedToAccount.Balance)
}

func TestStore_TransferTxDeadLock(t *testing.T) {
	s := db.NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// test concurrent transfer transactions
	n := 10
	amount := util.RandomAmount()
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := s.TransferTx(context.Background(), db.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final accounts
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	assert.Equal(t, account1.Balance, updatedFromAccount.Balance)
	assert.Equal(t, account2.Balance, updatedToAccount.Balance)
}
