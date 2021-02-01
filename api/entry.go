package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetEntryByIDRequest holds the params for get entry handler.
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
