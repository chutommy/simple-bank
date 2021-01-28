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

		// TODO:
		//   - check fromaccount
		//   - check toaccount
	}
}
