package expenses

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	t.Run("get all expenses succesfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		time, err := time.Parse("2006-01-02", "2024-05-18")

		if err != nil {
			t.Fatal(err)
		}

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, time, 1000, "food", "expense", "note", "image_url", 1).
			AddRow(2, time, 2000, "food", "expense", "note", "image_url", 2).
			AddRow(3, time, 3000, "food", "expense", "note", "image_url", 1)

		mock.ExpectPrepare(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense'`).WillBeClosed().ExpectQuery().WithoutArgs().WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err = h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		fmt.Println(rec.Body.String())
		assert.JSONEq(t, `[
			{"amount":1000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":1, "image_url":"image_url", "note":"note", "spender_id":1, "transaction_type":"expense"},
			{"amount":2000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":2, "image_url":"image_url", "note":"note", "spender_id":2, "transaction_type":"expense"},
			{"amount":3000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":3, "image_url":"image_url", "note":"note", "spender_id":1, "transaction_type":"expense"}
	]`, rec.Body.String())
	})

	t.Run("get all expenses failed on scan", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, "invalid date", 1000, "food", "expense", "note", "image_url", 1)

		mock.ExpectQuery(`SELECT id,date,amount,category,transaction_type,note,image_url,spender_id FROM transaction WHERE transaction_type = 'expense'`).WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get all expenses failed on database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

		mock.ExpectQuery(`SELECT id,date,amount,category,transaction_type,note,image_url,spender_id FROM transaction WHERE transaction_type = 'expense'`).WillReturnError(assert.AnError)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("get all expenses succesfully with empty result", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"})

		mock.ExpectPrepare(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense'`).ExpectQuery().WithoutArgs().WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[]`, rec.Body.String())
	})

	t.Run("get expenses with date query params should call query with date arg", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.QueryParams().Add("date", "2024-05-18")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		time1, _ := time.Parse("2006-01-02", "2024-05-18")

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, time1, 1000, "food", "expense", "note", "image_url", 1)

		mock.ExpectPrepare(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense' AND date = $1`).ExpectQuery().WithArgs("2024-05-18").WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[
			{"amount":1000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":1, "image_url":"image_url", "note":"note", "spender_id":1, "transaction_type":"expense"}
	]`, rec.Body.String())
	})
	t.Run("get expenses with amount query params should call query with amount arg", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.QueryParams().Add("amount", "1000")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		time1, _ := time.Parse("2006-01-02", "2024-05-18")

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, time1, 1000, "food", "expense", "note", "image_url", 1)

		mock.ExpectPrepare(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense' AND amount = $1`).ExpectQuery().WithArgs("1000").WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[
			{"amount":1000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":1, "image_url":"image_url", "note":"note", "spender_id":1, "transaction_type":"expense"}
	]`, rec.Body.String())
	})
	t.Run("get expenses with category query params should call query with category arg", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.QueryParams().Add("category", "food")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		time1, _ := time.Parse("2006-01-02", "2024-05-18")

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, time1, 1000, "food", "expense", "note", "image_url", 1)

		mock.ExpectPrepare(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense' AND category = $1`).ExpectQuery().WithArgs("food").WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetExpenses(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[
			{"amount":1000, "category":"food", "date":"2024-05-18T00:00:00Z", "id":1, "image_url":"image_url", "note":"note", "spender_id":1, "transaction_type":"expense"}
	]`, rec.Body.String())
	})

}
