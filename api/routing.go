package api

import (
	"github.com/gin-gonic/gin"
)

// getRouter sets up the routing for the given Server and returns
// the constructed gin router.
func getRouter(s *Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	accounts := r.Group("/accounts")
	{
		accounts.POST("", s.createAccount)
		accounts.GET("/:id", s.getAccountByID)
		accounts.GET("", s.listAccounts)
		accounts.PUT("/:id", s.updateAccount)
		accounts.DELETE("/:id", s.deleteAccount)
	}

	entries := r.Group("/entries")
	{
		entries.GET("/:id", s.getEntryByID)
	}

	return r
}
