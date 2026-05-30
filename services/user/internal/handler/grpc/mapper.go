package grpc

import (
	"user-service/internal/domain/entities"
	"user-service/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoUser(user *entities.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
