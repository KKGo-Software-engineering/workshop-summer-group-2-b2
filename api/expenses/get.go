package expenses

import (
	"database/sql"
	"net/http"

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

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT id,date,amount,category,transaction_type,note,image_url,spender_id FROM transaction WHERE transaction_type = 'expense'`)
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
