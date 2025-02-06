-- name: GetAccountPasswordByID :one
SELECT id, password_hash FROM accounts WHERE username = $1;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE username = $1;

-- name: UpdateAccountPassword :exec
UPDATE accounts SET password_hash = $1 WHERE id = $2;

-- name: GetAccountByUsername :one
SELECT * FROM accounts WHERE username = $1;