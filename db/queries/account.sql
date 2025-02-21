-- name: GetAccountPasswordHashByUsername :one
SELECT id, password_hash
FROM accounts
WHERE username = $1;

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: UpdateAccountPassword :exec
UPDATE accounts
SET password_hash = $1
WHERE id = $2;

-- name: GetAccountByUsername :one
SELECT *
FROM accounts
WHERE username = $1;


-- name: CreateAccount :one
INSERT INTO accounts (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING *;


-- sqlc.arg(password_hash) json:"-"
