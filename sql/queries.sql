-- name: CreateAccount :exec
INSERT INTO account (account_state, balance, client_name)
VALUES (true, ?, ?);

-- name: GetAllAcounts :many
SELECT * FROM account;

-- name: GetAcountById :one
SELECT * FROM account
WHERE account_number = ?;

-- name: DeleteAcount :exec
UPDATE account
SET account_state = false
WHERE account_number = ?;

-- name: ModifyBalance :exec
UPDATE account
SET balance = balance + sqlc.arg(balance)
WHERE account_number = ?;