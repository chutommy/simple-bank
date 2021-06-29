package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chutommy/simple-bank/laptop/laptop"
	"github.com/chutommy/simple-bank/laptop/pkg/repo"
)

type LaptopServer struct {
	Repo repo.Repo
}

func NewLaptopServer(repo repo.Repo) *LaptopServer {
	return &LaptopServer{
		Repo: repo,
	}
}

func (l *LaptopServer) CreateLaptop(ctx context.Context, req *laptop.CreateLaptopRequest) (*laptop.CreateLaptopResponse, error) {
	newLaptop := req.GetLaptop()

	// process laptop's ID
	if newLaptop.Id != "" {
		_, err := uuid.Parse(newLaptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %w", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not generate a new random UUID: %w", err)
		}

		newLaptop.Id = id.String()
	}

	// store laptop
	err := l.Repo.CreateLaptop(newLaptop)
	if err != nil {
		code := codes.Internal
		switch {
		case errors.Is(err, repo.ErrAlreadyExists):
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop: %w", err)
	}

	// response
	resp := &laptop.CreateLaptopResponse{
		Id: newLaptop.Id,
	}

	return resp, nil
}

func (l *LaptopServer) mustEmbedUnimplementedLaptopServiceServer() {
	panic("implement me")
}
