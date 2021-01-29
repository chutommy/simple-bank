package api

import (
	"github.com/gin-gonic/gin"
)

func getRouter(s *Server) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/accounts", s.createAccount)

	return r
}
