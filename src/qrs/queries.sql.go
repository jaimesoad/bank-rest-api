// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package qrs

import (
	"context"
)

const createAccount = `-- name: CreateAccount :exec
INSERT INTO account (account_state, balance, client_name)
VALUES (true, ?, ?)
`

type CreateAccountParams struct {
	Balance    float64 `json:"balance"`
	ClientName string  `json:"clientName"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	_, err := q.db.ExecContext(ctx, createAccount, arg.Balance, arg.ClientName)
	return err
}

const deleteAcount = `-- name: DeleteAcount :exec
UPDATE account
SET account_state = false
WHERE account_number = ?
`

func (q *Queries) DeleteAcount(ctx context.Context, accountNumber int32) error {
	_, err := q.db.ExecContext(ctx, deleteAcount, accountNumber)
	return err
}

const getAcountById = `-- name: GetAcountById :one
SELECT account_number, account_state, balance, client_name FROM account
WHERE account_number = ?
`

func (q *Queries) GetAcountById(ctx context.Context, accountNumber int32) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAcountById, accountNumber)
	var i Account
	err := row.Scan(
		&i.AccountNumber,
		&i.AccountState,
		&i.Balance,
		&i.ClientName,
	)
	return i, err
}

const getAllAcounts = `-- name: GetAllAcounts :many
SELECT account_number, account_state, balance, client_name FROM account
`

func (q *Queries) GetAllAcounts(ctx context.Context) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, getAllAcounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.AccountNumber,
			&i.AccountState,
			&i.Balance,
			&i.ClientName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const modifyBalance = `-- name: ModifyBalance :exec
UPDATE account
SET balance = balance + ?
WHERE account_number = ?
`

type ModifyBalanceParams struct {
	Balance       float64 `json:"balance"`
	AccountNumber int32   `json:"accountNumber"`
}

func (q *Queries) ModifyBalance(ctx context.Context, arg ModifyBalanceParams) error {
	_, err := q.db.ExecContext(ctx, modifyBalance, arg.Balance, arg.AccountNumber)
	return err
}
