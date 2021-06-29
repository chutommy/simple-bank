package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/chutommy/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidAddress is returned when an invalid server address is provided.
	ErrInvalidAddress = errors.New("invalid server address")
	// ErrServerShutdownTimeout is returned when the process of shutting down ss too long.
	ErrServerShutdownTimeout = errors.New("shutdown timeout, server forced to shutdown")
)

// Server serves HTTP requests for the banking service.
type Server struct {
	store  db.Store
	router *gin.Engine

	Srv *http.Server
}

// NewServer constructs a new HTTP Server and setup the routing.
func NewServer(store db.Store) *Server {
	s := &Server{store: store}
	s.router = getRouter(s)

	s.Srv = &http.Server{
		Handler: s.router,
	}

	return s
}

// Start initializes the HTTP Server on a given address.
func (s *Server) Start(addr string) error {
	// get address's port number
	addrSplit := strings.Split(addr, ":")
	if len(addrSplit) != 2 {
		return ErrInvalidAddress
	}

	s.Srv.Addr = fmt.Sprintf(":%s", addrSplit[1])

	return s.Srv.ListenAndServe()
}

// Stop gracefully shutdowns the server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Srv.Shutdown(ctx); err != nil {
		return ErrServerShutdownTimeout
	}

	return nil
}

// errorResponse serializes the error into a JSON key-value pair struct.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
