package api

import (
	"database/sql"
	"errors"
	db "github.com/chutommy/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MakeTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=0"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=0"`
	Amount        int64 `json:"amount" binding:"required"`
}

func (s *Server) makeTransfer(c *gin.Context) {
	var req MakeTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	result, err := s.store.TransferTx(c, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, result)
}
