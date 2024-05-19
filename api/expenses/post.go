package expenses

import (
	"net/http"

	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/KKGo-Software-engineering/workshop-summer/api/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = `INSERT INTO transaction (date, amount, category, transaction_type, note, image_url, spender_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
)

func (h handler) Create(c echo.Context) error {
	logger := mlog.L(c)
	var transaction model.Transaction
	ctx := c.Request().Context()
	err := c.Bind(&transaction)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt,
		transaction.Date,
		transaction.Amount,
		transaction.Category,
		transaction.TransactionType,
		transaction.Note,
		transaction.ImageURL,
		transaction.SpenderId,
	).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	transaction.ID = lastInsertId
	return c.JSON(http.StatusCreated, transaction)
}
