package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateAccountRequest holds parameters for createAccount handler.
type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required,ascii"`
	Currency string `json:"currency" binding:"required,uppercase"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	// store the new account into the database
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

// GetAccountByIDRequest holds parameters for getAccountByID handler.
type GetAccountByIDRequest struct {
	ID int64 `uri:"id" binding:"required,numeric,min=1"`
}

func (s *Server) getAccountByID(c *gin.Context) {
	var req GetAccountByIDRequest
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

// ListAccountsRequest holds parameters for listAccounts handler.
type ListAccountsRequest struct {
	PageNum  int32 `form:"page_num" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=1000"`
}

func (s *Server) listAccounts(c *gin.Context) {
	var req ListAccountsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	// query accounts
	params := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageNum - 1) * req.PageSize,
	}

	accounts, err := s.store.ListAccounts(c, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, accounts)
}

// UpdateAccountRequestURI holds URI parameters for updateAccount handler.
type UpdateAccountRequestURI struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// UpdateAccountRequestJSON holds JSON parameters for updateAccount handler.
type UpdateAccountRequestJSON struct {
	Balance int64 `json:"balance" binding:"required,numeric"`
}

func (s *Server) updateAccount(c *gin.Context) {
	var reqURI UpdateAccountRequestURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var reqJSON UpdateAccountRequestJSON
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	account, err := s.store.UpdateAccountBalance(c, db.UpdateAccountBalanceParams{
		ID:      reqURI.ID,
		Balance: reqJSON.Balance,
	})
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

// DeleteAccountRequest holds URI params to delete an db.Account.
type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteAccount(c *gin.Context) {
	var req DeleteAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if err := s.store.DeleteAccount(c, req.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}
	}

	c.JSON(http.StatusOK, nil)
}
