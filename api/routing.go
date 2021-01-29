package api

import (
	"github.com/gin-gonic/gin"
)

// getRouter sets up the routing for the given Server and returns
// the constructed gin router.
func getRouter(s *Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/accounts", s.createAccount)
	r.GET("/accounts/:id", s.getAccountByID)

	return r
}