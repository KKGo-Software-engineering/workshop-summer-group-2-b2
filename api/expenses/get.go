package expenses

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/KKGo-Software-engineering/workshop-summer/api/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

func (h handler) GetExpenses(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	params := c.QueryParams()

	date := params.Get("date")
	amount := params.Get("amount")
	category := params.Get("category")
	page := params.Get("page")

	pageSize := 1

	query := "SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction WHERE transaction_type = 'expense'"
	var conditions []string
	var args []interface{}

	if date != "" {
		conditions = append(conditions, "date = $1")
		args = append(args, date)
	}

	if amount != "" {
		conditions = append(conditions, fmt.Sprintf("amount = $%d", len(args)+1))
		args = append(args, amount)
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", len(args)+1))
		args = append(args, category)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	var pageNumber int
	var err error

	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid page number")
		}

		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, (pageNumber-1)*pageSize)
	}

	stmt, err := h.db.PrepareContext(ctx, query)

	if err != nil {
		logger.Error("prepare error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	rows, err := stmt.QueryContext(ctx, args...)

	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()

	var txs []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		err := rows.Scan(&tx.ID, &tx.Date, &tx.Amount, &tx.Category, &tx.TransactionType, &tx.Note, &tx.ImageURL, &tx.SpenderId)

		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		txs = append(txs, tx)
	}

	if len(txs) == 0 {
		return c.JSON(http.StatusOK, []model.Transaction{})
	}

	return c.JSON(http.StatusOK, txs)
}
