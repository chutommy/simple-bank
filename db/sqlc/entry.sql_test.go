package db_test

import (
	"context"
	"database/sql"
	"testing"

	db "github.com/chutommy/simple-bank/db/sqlc"
	"github.com/chutommy/simple-bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account db.Account) db.Entry {
	t.Helper()

	arg := db.CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(),
	}

	// create a new entry
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	// check values
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	acc1 := createRandomAccount(t)
	createRandomEntry(t, acc1)
}

func TestQueries_GetEntry(t *testing.T) {
	acc1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc1)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	if assert.NotEmpty(t, entry2) {
		assert.Equal(t, entry1.ID, entry2.ID)
		assert.Equal(t, entry1.AccountID, entry2.AccountID)
		assert.Equal(t, entry1.Amount, entry2.Amount)
		assert.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	}
}

func TestQueries_ListEntries(t *testing.T) {
	acc1 := createRandomAccount(t)

	// generate entries
	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc1)
	}

	arg := db.ListEntriesParams{
		AccountID: acc1.ID,
		Limit:     5,
		Offset:    5,
	}

	// retrieve list
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, int(arg.Limit))

	// check entries
	for _, entry := range entries {
		if assert.NotEmpty(t, entry) {
			assert.Equal(t, acc1.ID, entry.AccountID)
		}
	}
}

func TestQueries_UpdateEntryAmount(t *testing.T) {
	acc1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc1)

	arg := db.UpdateEntryAmountParams{
		ID:     entry1.ID,
		Amount: util.RandomAmount(),
	}

	// update entry
	entry2, err := testQueries.UpdateEntryAmount(context.Background(), arg)
	require.NoError(t, err)

	// check values
	if assert.NotEmpty(t, entry2) {
		assert.Equal(t, entry1.ID, entry2.ID)
		assert.Equal(t, entry1.AccountID, entry2.AccountID)
		assert.Equal(t, arg.Amount, entry2.Amount)
		assert.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	}
}

func TestQueries_DeleteEntry(t *testing.T) {
	acc1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc1)

	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)

	// check the deleted object does not exist anymore
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, entry2)
}
