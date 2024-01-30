package main

import (
	"bank/src/qrs"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

//go:embed sql/schema.sql
var schema string

func main() {
	q, ctx, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	e := echo.New()

	r := e.Group("/api/v1/accounts")

	r.POST("", func(c echo.Context) error {
		var account qrs.CreateAccountParams

		err := c.Bind(&account)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if account.ClientName == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		err = q.CreateAccount(ctx, account)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusCreated, "Account created succesfully")
	})

	r.GET("", func(c echo.Context) error {

		accounts, err := q.GetAllAcounts(ctx)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if len(accounts) == 0 {
			accounts = []qrs.Account{}
		}

		return c.JSON(http.StatusOK, accounts)
	})

	r.GET("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		account, err := q.GetAcountById(ctx, int32(id))
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, account)
	})

	r.DELETE("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		err = q.DeleteAcount(ctx, int32(id))
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, "Account deleted successfully")
	})

	r.PUT("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		account, err := q.GetAcountById(ctx, int32(id))
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	
		if !account.AccountState {
			return c.String(http.StatusBadRequest, "Account already deleted")
		}

		var newBalance = struct {
			Balace float64 `json:"balance"`
		}{}

		err = c.Bind(&newBalance)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		var balance = qrs.ModifyBalanceParams {
			Balance: newBalance.Balace,
			AccountNumber: int32(id),
		}

		err = q.ModifyBalance(ctx, balance)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, "Balance modified successfully")
	})

	log.Fatal(e.Start(":3000"))
}

func connectDB() (*qrs.Queries, context.Context, error) {
	conf := mysql.Config{
		User:   "user",
		Passwd: "passwd",
		DBName: "bank",
		Addr:   "localhost",
	}

	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, nil, err
	}

	q := qrs.New(db)

	return q, context.Background(), nil
}
