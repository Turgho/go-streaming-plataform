package usecase

import (
	"context"
	"fmt"

	"user-service/internal/domain/entities"
	domainrepo "user-service/internal/domain/repositories"
	"user-service/internal/infra/jwt"
	"user-service/pkg/hash"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUseCase struct {
	repo domainrepo.UserRepository
}

func NewUserUseCase(repo domainrepo.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Register(ctx context.Context, username, email, password, confirmPassword string) (*entities.User, error) {
	if password != confirmPassword {
		return nil, status.Error(codes.InvalidArgument, "senhas não iguais")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Error(codes.Internal, "gerar id do usuário")
	}

	passwordHash := hash.HashPassword(password)

	user, err := entities.NewUser(id.String(), username, email, passwordHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (token string, user *entities.User, err error) {
	if email == "" || password == "" {
		return "", nil, status.Error(codes.InvalidArgument, "email ou senha inválidos")
	}

	user, err = uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, status.Error(codes.InvalidArgument, "credenciais inválidas")
	}

	if !hash.VerifyPassword(password, user.PasswordHash) {
		return "", nil, status.Error(codes.InvalidArgument, "credenciais inválidas")
	}

	token, err = jwt.Generate(user.ID, user.Email)
	if err != nil {
		return "", nil, status.Error(codes.Internal, "gerar token")
	}

	return token, user, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*entities.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return user, nil
}

func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, mapRepoError(err)
	}
	return user, nil
}

func (uc *UserUseCase) ValidateToken(ctx context.Context, token string) (userID, email string, err error) {
	claims, err := jwt.Validate(token)
	if err != nil {
		return "", "", status.Errorf(codes.Unauthenticated, "token inválido: %v", err)
	}

	user, err := uc.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", status.Errorf(codes.NotFound, "usuário não encontrado")
	}

	return user.ID, user.Email, nil
}

func mapRepoError(err error) error {
	if err == nil {
		return nil
	}
	return status.Error(codes.NotFound, fmt.Sprintf("usuário não encontrado: %v", err))
}
