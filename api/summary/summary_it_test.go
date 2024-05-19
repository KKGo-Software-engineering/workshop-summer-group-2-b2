//go:build integration

package summary

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/migration"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestExpenseSummaryIT(t *testing.T) {
	t.Run("get expense summary successfully", func(t *testing.T) {
		sql := newDatabase(t)

		h := New(sql)
		e := echo.New()
		defer e.Close()

		e.GET("/expenses/summary", h.GetSummary)

		req := httptest.NewRequest(http.MethodGet, "/expenses/summary", nil)
		q := req.URL.Query()
		q.Add("spender_id", "1")
		req.URL.RawQuery = q.Encode()
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"total_income": 2000,
			"total_expenses": 500,
			"current_balance": 1500
		}`, rec.Body.String())
	})
}

func newDatabase(t *testing.T) *sql.DB {
	cfg := config.Parse("DOCKER")
	sql, err := sql.Open("postgres", cfg.PostgresURI())
	if err != nil {
		t.Fatal(err)
	}
	migration.ApplyMigrations(sql)
	sql.Query("INSERT INTO public.transaction (amount, transaction_type, spender_id) VALUES (1000, 'income',1);")
	sql.Query("INSERT INTO public.transaction (amount, transaction_type, spender_id) VALUES (500, 'expense',1);")
	sql.Query("INSERT INTO public.transaction (amount, transaction_type, spender_id) VALUES (1000, 'income',1);")
	t.Cleanup(func() {
		sql.Query("DELETE FROM public.transaction WHERE spender_id = 1;")
		sql.Close()
	})
	return sql
}
