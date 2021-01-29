package api

import (
	db "github.com/chutified/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for the banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer constructs a new HTTP Server and setup the routing.
func NewServer(store *db.Store) *Server {
	s := &Server{store: store}
	s.router = getRouter(s)

	return s
}

// Start initializes the HTTP Server on a given address.
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// errorResponse serializes the error into a JSON key-value pair struct.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
