package server

import (
	"context"
	"user-service/internal/domain/entities"
	domain "user-service/internal/domain/repositories"
	"user-service/internal/infra/jwt"
	"user-service/pkg/hash"
	"user-service/pkg/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	repo domain.UserRepository
}

func NewUserServer(repo domain.UserRepository) *Server {
	return &Server{repo: repo}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate user id")
	}

	if req.Password != req.ConfirmPassword {
		return nil, status.Error(codes.InvalidArgument, "senhas não iguais")
	}

	passwordHash := hash.HashPassword(req.Password)

	user, err := entities.NewUser(id.String(), req.Username, req.Email, passwordHash)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return toProtoUser(user), nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email ou senha inválidos")
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !hash.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, status.Error(codes.InvalidArgument, "senha inválidos")
	}

	token, err := jwt.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
		User:  toProtoUser(user),
	}, nil
}

func (s *Server) GetByEmail(ctx context.Context, req *pb.GetByEmailRequest) (*pb.User, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return toProtoUser(user), nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.User, error) {
	user, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return toProtoUser(user), nil
}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := jwt.Validate(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token inválido: %v", err)
	}

	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "usuário não encontrado: %v", err)
	}

	return &pb.ValidateTokenResponse{
		UserId: user.ID,
		Email:  user.Email,
	}, nil
}

// Helpers
func toProtoUser(user *entities.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
