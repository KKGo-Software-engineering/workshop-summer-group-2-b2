package summary

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetSummary(t *testing.T) {
	t.Run("get spender expenses summary", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		q := req.URL.Query()
		q.Add("spender_id", "1")
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"total_income", "total_expense", "current_balance"}).
			AddRow(1000, 500, 500)
		mock.ExpectQuery(` WITH balance_summary AS (
			SELECT
				SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income,
				SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses
			FROM public."transaction"
			WHERE spender_id = $1 -- Assuming you want to filter by a specific spender
		)
		SELECT
			total_income,
			total_expenses,
			total_income - total_expenses AS current_balance
		FROM balance_summary;`).WillReturnRows(rows)
		h := New(db)
		err := h.GetSummary(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"total_income": 1000,
			"total_expenses": 500,
			"current_balance": 500
		}`, rec.Body.String())
	})

	t.Run("get expense summar failed on database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(` WITH balance_summary AS (
			SELECT
				SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income,
				SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses
			FROM public."transaction"
			WHERE spender_id = $1 -- Assuming you want to filter by a specific spender
		)
		SELECT
			total_income,
			total_expenses,
			total_income - total_expenses AS current_balance
		FROM balance_summary;`).WillReturnError(assert.AnError)

		h := New(db)
		err := h.GetSummary(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
