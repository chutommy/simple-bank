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

func createRandomUser(t *testing.T) db.User {
	t.Helper()

	// construct params
	arg := db.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FirstName:      util.RandomOwner(),
		LastName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	// create account
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	// check values
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user1 := createRandomUser(t)

	// get account
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)

	// compare values
	if assert.NotEmpty(t, user2) {
		assert.Equal(t, user1.Username, user2.Username)
		assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
		assert.Equal(t, user1.FirstName, user2.FirstName)
		assert.Equal(t, user1.LastName, user2.LastName)
		assert.Equal(t, user1.Email, user2.Email)
	}
}

func TestQueries_UpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)

	arg := db.UpdateUserPasswordParams{
		Username:       user1.Username,
		HashedPassword: "new_password",
	}

	// update password
	user2, err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)

	// compare value
	if assert.NotEmpty(t, user2) {
		assert.Equal(t, user1.Username, user2.Username)
		assert.Equal(t, "new_password", user2.HashedPassword)
		assert.Equal(t, user1.FirstName, user2.FirstName)
		assert.Equal(t, user1.LastName, user2.LastName)
		assert.Equal(t, user1.Email, user2.Email)
	}
}

func TestQueries_DeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	// delete user
	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)

	// check
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Empty(t, user2)
}
