package api

import (
	"database/sql"
	"errors"
	db "github.com/chutified/simple-bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetEntryByIDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getEntryByID(c *gin.Context) {
	var req GetEntryByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	entry, err := s.store.GetEntry(c, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, entry)
}

type ListEntriesRequestURI struct {
	AccountID int64 `uri:"account_id" binding:"required,min=1"`
}

type ListEntriesRequestQuery struct {
	PageNum  int32 `form:"page_num" binding:"required,numeric,min=1"`
	PageSize int32 `form:"page_size" binding:"required,numeric,min=1"`
}

func (s *Server) listEntries(c *gin.Context) {
	var reqURI ListEntriesRequestURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var reqQuery ListEntriesRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	entries, err := s.store.ListEntries(c, db.ListEntriesParams{
		AccountID: reqURI.AccountID,
		Limit:     reqQuery.PageSize,
		Offset:    (reqQuery.PageNum - 1) * reqQuery.PageSize,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, entries)
}

type CreateEntryRequest struct {
	AccountID int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required,min=1"`
}

func (s *Server) createEntry(c *gin.Context) {
	var req CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	entry, err := s.store.CreateEntry(c, db.CreateEntryParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, entry)
}

type UpdateEntryRequestURI struct {
	ID int64 `uri:"id" binding:"required,min=0"`
}

type UpdateEntryRequestJSON struct {
	Amount int64 `json:"amount" binding:"required"`
}

func (s *Server) updateEntry(c *gin.Context) {
	var reqURI UpdateEntryRequestURI
	if err := c.ShouldBindUri(&reqURI); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var reqJSON UpdateEntryRequestJSON
	if err := c.ShouldBindJSON(&reqJSON); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	entry, err := s.store.UpdateEntryAmount(c, db.UpdateEntryAmountParams{
		ID:     reqURI.ID,
		Amount: reqJSON.Amount,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, entry)
}

type DeleteEntryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteEntry(c *gin.Context) {
	var req DeleteEntryRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if err := s.store.DeleteEntry(c, req.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
		}

		return
	}

	c.JSON(http.StatusOK, nil)
}
