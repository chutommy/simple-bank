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

func createRandomTransfer(t *testing.T, from, to db.Account) db.Transfer {
	t.Helper()

	arg := db.CreateTransferParams{
		FromAccountID: from.ID,
		ToAccountID:   to.ID,
		Amount:        util.RandomAmount(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestQueries_GetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	if assert.NotEmpty(t, transfer2) {
		assert.Equal(t, transfer1.ID, transfer2.ID)
		assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
		assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
		assert.Equal(t, transfer1.Amount, transfer2.Amount)
		assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	}
}

func TestQueries_ListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}

	arg := db.ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, int(arg.Limit))

	for _, transfer := range transfers {
		if assert.NotEmpty(t, transfer) {
			assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
			assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
		}
	}
}

func TestQueries_UpdateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	arg := db.UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomAmount(),
	}

	transfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)

	if assert.NotEmpty(t, transfer2) {
		assert.Equal(t, transfer1.ID, transfer2.ID)
		assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
		assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
		assert.Equal(t, arg.Amount, transfer2.Amount)
		assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt)
	}
}

func TestQueries_DeleteTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, transfer2)
}
