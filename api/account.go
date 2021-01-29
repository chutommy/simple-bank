package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required,ascii"`
	Currency string `json:"currency" binding:"required,uppercase"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	params := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	c.JSON(http.StatusOK, account)
}

type getAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,numeric,min=1"`
}

func (s *Server) getAccountByID(c *gin.Context) {
	var req getAccountByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	// query the account ID
	account, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, account)
}
