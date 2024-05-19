package summary

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	db *sql.DB
}

const (
	// exprenseStmt = `SELECT amount FROM transaction WHERE spender_id = $1 AND transaction_type = 'expense'`

	summaryStmt = ` WITH balance_summary AS (
        SELECT
            SUM(CASE WHEN transaction_type = 'income' THEN amount ELSE 0 END) AS total_income,
            SUM(CASE WHEN transaction_type = 'expense' THEN amount ELSE 0 END) AS total_expenses
        FROM public.transaction
        WHERE spender_id = $1 -- Assuming you want to filter by a specific spender
    )
    SELECT
        total_income,
        total_expenses,
        total_income - total_expenses AS current_balance
    FROM balance_summary;`
)

func New(db *sql.DB) *handler {
	return &handler{db}
}

type SummaryResponse struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpenses  float64 `json:"total_expenses"`
	CurrentBalance float64 `json:"current_balance"`
}

func (h handler) GetSummary(c echo.Context) error {
	var query struct {
		ID int `query:"spender_id"`
	}
	c.Bind(&query)
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, summaryStmt, query.ID)

	if err != nil {
		fmt.Println("🚀 | file: summary.go | line 54 | func | err : ", err)
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var responseObject = SummaryResponse{}
	for rows.Next() {
		err := rows.Scan(&responseObject.TotalIncome, &responseObject.TotalExpenses, &responseObject.CurrentBalance)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, responseObject)
}
