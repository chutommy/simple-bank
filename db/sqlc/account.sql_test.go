package db_test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/chutified/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) db.Account {
	t.Helper()

	user := createRandomUser(t)

	// construct params
	arg := db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	// create account
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	// check values
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestQueries_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestQueries_GetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	// get account
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	// compare values
	if assert.NotEmpty(t, acc2) {
		assert.Equal(t, acc1.ID, acc2.ID)
		assert.Equal(t, acc1.Owner, acc2.Owner)
		assert.Equal(t, acc1.Balance, acc2.Balance)
		assert.Equal(t, acc1.Currency, acc2.Currency)
		assert.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	}
}

func TestQueries_GetAccountForUpdate(t *testing.T) {
	acc1 := createRandomAccount(t)

	// get account
	acc2, err := testQueries.GetAccountForUpdate(context.Background(), acc1.ID)
	require.NoError(t, err)

	// compare values
	if assert.NotEmpty(t, acc2) {
		assert.Equal(t, acc1.ID, acc2.ID)
		assert.Equal(t, acc1.Owner, acc2.Owner)
		assert.Equal(t, acc1.Balance, acc2.Balance)
		assert.Equal(t, acc1.Currency, acc2.Currency)
		assert.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	}
}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := db.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	// list accounts
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	assert.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}
}

func TestQueries_UpdateAccountBalance(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := db.UpdateAccountBalanceParams{
		ID:      acc1.ID,
		Balance: util.RandomBalance(),
	}

	// update balance
	acc2, err := testQueries.UpdateAccountBalance(context.Background(), arg)
	require.NoError(t, err)

	// compare value
	if assert.NotEmpty(t, acc2) {
		assert.Equal(t, acc1.ID, acc2.ID)
		assert.Equal(t, acc1.Owner, acc2.Owner)
		assert.Equal(t, arg.Balance, acc2.Balance)
		assert.Equal(t, acc1.Currency, acc2.Currency)
		assert.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	}
}

func TestQueries_AddAccountBalance(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := db.AddAccountBalanceParams{
		ID:     acc1.ID,
		Amount: util.RandomBalance(),
	}

	// update balance
	acc2, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)

	// compare value
	if assert.NotEmpty(t, acc2) {
		assert.Equal(t, acc1.ID, acc2.ID)
		assert.Equal(t, acc1.Owner, acc2.Owner)
		assert.Equal(t, acc1.Balance+arg.Amount, acc2.Balance)
		assert.Equal(t, acc1.Currency, acc2.Currency)
		assert.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	}
}

func TestQueries_DeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	// delete account
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	// check
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, acc2)
}
