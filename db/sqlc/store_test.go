package db_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

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

		// TODO:
		//   - check transfer
		//   - check fromentry
		//   - check toentry
		//   - check fromaccount
		//   - check toaccount
	}
}
