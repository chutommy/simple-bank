// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, hashed_password, first_name, last_name, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING username, hashed_password, first_name, last_name, email, password_modified_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FirstName,
		arg.LastName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordModifiedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE username = $1
`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, username)
	return err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, first_name, last_name, email, password_modified_at, created_at
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordModifiedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE username = $1
RETURNING username, hashed_password, first_name, last_name, email, password_modified_at, created_at
`

type UpdateUserPasswordParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.Username, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordModifiedAt,
		&i.CreatedAt,
	)
	return i, err
}