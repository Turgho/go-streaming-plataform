package interceptor

import (
	"context"
	"strings"
	"user-service/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(userClient pb.UserServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata ausente")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token ausente")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")

		user, err := userClient.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "token inválido")
		}

		// injeta o usuário no contexto
		ctx = context.WithValue(ctx, "user", user)
		return handler(ctx, req)
	}
}
