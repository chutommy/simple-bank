package api

import (
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

func (s *Server) createUser(c *gin.Context) {}

type GetUserRequest struct {
	Username string
}

func (s *Server) getUser(c *gin.Context) {}

type UpdateUserPasswordRequest struct {
	Username     string
	OldPassoword string
	NewPassword  string
}

func (s *Server) updateUserPassword(c *gin.Context) {}

type DeleteUserRequest struct {
	Username string
	Password string
}

func (s *Server) deleteUser(c *gin.Context) {}
