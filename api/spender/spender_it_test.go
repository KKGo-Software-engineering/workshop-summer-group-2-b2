//go:build integration

package spender

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/migration"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)
		payload := `{"name": "HongJot", "email": "hong@jot.ok"}`
		req := httptest.NewRequest(http.MethodPost, "/spenders", strings.NewReader(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})
}

func TestGetAllSpenderIT(t *testing.T) {
	t.Run("get all spender successfully", func(t *testing.T) {
		sql := newDatabase(t)

		h := New(config.FeatureFlag{}, sql)
		e := echo.New()
		defer e.Close()

		e.GET("/spenders", h.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/spenders", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
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
		sql.Close()
	})
	return sql
}
