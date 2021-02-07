-- name: CreateUser :one
INSERT INTO users (username, hashed_password, first_name, last_name, email, password_modified_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $3
WHERE username = $1
  AND hashed_password = $2
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE username = $1
  AND hashed_password = $2;
