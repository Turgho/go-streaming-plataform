package interceptor

import (
	"context"
	"log"
	"strings"

	userpb "upload-service/pkg/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(userClient userpb.UserServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log.Println("AuthInterceptor chamado")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata ausente")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token ausente")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")

		user, err := userClient.ValidateToken(ctx, &userpb.ValidateTokenRequest{Token: token})
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token inválido")
		}

		ctx = context.WithValue(ctx, "user_id", user.UserId)
		return handler(ctx, req)
	}
}
