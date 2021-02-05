package api

import (
	"github.com/gin-gonic/gin"
)

type MakeTransferRequest struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

func (s *Server) makeTransfer(c *gin.Context) {
}
