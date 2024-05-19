//go:build integration

package expenses

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/model"
	"github.com/KKGo-Software-engineering/workshop-summer/migration"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// func getTestDatabaseFromConfig() (*sql.DB, error) {
// 	cfg := config.Parse("DOCKER")
// 	sql, err := sql.Open("postgres", cfg.PostgresURI())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sql, nil
// }

func TestCreateExpensesIT(t *testing.T) {
	t.Run("create Expense successfully when feature toggle is enable", func(t *testing.T) {
		transactions := model.Transaction{
			Date:            time.Date(2024, 5, 19, 20, 10, 0, 0, time.UTC),
			Amount:          999.75,
			Category:        "Groceries",
			TransactionType: "Expense",
			Note:            "Weekly grocery shopping",
			ImageURL:        "http://example.com/receipt.jpg",
			SpenderId:       5,
		}
		sql := newDatabase(t)

		h := New(config.FeatureFlag{EnableCreateSpender: true}, sql)
		e := echo.New()
		defer e.Close()

		e.POST("/expenses", h.Create)

		jsonTransaction, _ := json.Marshal(transactions)
		req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewReader(jsonTransaction))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		var got model.Transaction
		err := json.Unmarshal(rec.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, transactions.Date, got.Date)
		assert.Equal(t, transactions.Amount, got.Amount)
		assert.Equal(t, transactions.Category, got.Category)
		assert.Equal(t, transactions.TransactionType, got.TransactionType)
		assert.Equal(t, transactions.Note, got.Note)
		assert.Equal(t, transactions.ImageURL, got.ImageURL)
		assert.Equal(t, transactions.SpenderId, got.SpenderId)
		// assert.JSONEq(t, `{"amount":999.75, "category":"Groceries", "date":"2024-05-19T20:10:00Z", "id":1, "image_url":"http://example.com/receipt.jpg", "note":"Weekly grocery shopping", "spender_id":5, "transaction_type":"Expense"}`, rec.Body.String())
	})
}

func newDatabase(t *testing.T) *sql.DB {
	cfg := config.Parse("DOCKER")
	sql, err := sql.Open("postgres", cfg.PostgresURI())
	if err != nil {
		t.Fatal(err)
	}
	migration.ApplyMigrations(sql)
	t.Cleanup(func() {
		sql.Query("DELETE FROM public.transaction WHERE amout = 999.75 ;")
		sql.Close()
	})
	return sql
}
