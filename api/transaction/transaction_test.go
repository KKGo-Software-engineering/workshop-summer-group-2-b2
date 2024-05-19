package transaction_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/KKGo-Software-engineering/workshop-summer/api/model"
	"github.com/KKGo-Software-engineering/workshop-summer/api/transaction"
)

func TestCreateSpender(t *testing.T) {

	t.Run("create spender succesfully when feature toggle is enable", func(t *testing.T) {
		e := echo.New()
		transactions := model.Transaction{
			Date:            time.Date(2024, 5, 19, 20, 10, 0, 0, time.UTC),
			Amount:          150.75,
			Category:        "Groceries",
			TransactionType: "Expense",
			Note:            "Weekly grocery shopping",
			ImageURL:        "http://example.com/receipt.jpg",
			SpenderId:       5,
		}
		defer e.Close()
		jsonTransaction, err := json.Marshal(transactions)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonTransaction))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(`INSERT INTO transaction (date, amount, category, transaction_type, note, image_url, spender_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`).
			WithArgs(transactions.Date, transactions.Amount, transactions.Category, transactions.TransactionType, transactions.Note, transactions.ImageURL, transactions.SpenderId).
			WillReturnRows(row)

		h := transaction.New(db)
		err = h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

}
