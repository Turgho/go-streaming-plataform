package grpc

import (
	"context"

	"user-service/internal/usecase"
	"user-service/pkg/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	users *usecase.UserUseCase
}

func NewServer(users *usecase.UserUseCase) *Server {
	return &Server{users: users}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.User, error) {
	user, err := s.users.Register(ctx, req.Username, req.Email, req.Password, req.ConfirmPassword)
	if err != nil {
		return nil, err
	}
	return toProtoUser(user), nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, user, err := s.users.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		Token: token,
		User:  toProtoUser(user),
	}, nil
}

func (s *Server) GetByEmail(ctx context.Context, req *pb.GetByEmailRequest) (*pb.User, error) {
	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return toProtoUser(user), nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.User, error) {
	user, err := s.users.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return toProtoUser(user), nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token obrigatório")
	}

	userID, email, err := s.users.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateTokenResponse{
		UserId: userID,
		Email:  email,
	}, nil
}
