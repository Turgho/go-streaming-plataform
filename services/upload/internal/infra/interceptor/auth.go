package interceptor

import (
	"context"
	"strings"

	userpb "upload-service/pkg/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const serviceKeyHeader = "x-service-key"

var internalMethods = map[string]bool{
	"/upload.UploadService/UpdateStatus": true,
}

// AuthInterceptor valida JWT de usuário ou chave interna para chamadas entre serviços.
func AuthInterceptor(userClient userpb.UserServiceClient, serviceKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// verificação de métodos para transcode-service.
		if internalMethods[info.FullMethod] {
			if !validateServiceKey(ctx, serviceKey) {
				return nil, status.Error(codes.Unauthenticated, "chave de serviço inválida")
			}
			return handler(ctx, req)
		}

		userID, err := authenticateUser(ctx, userClient)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, UserIDKey, userID)
		return handler(ctx, req)
	}
}

func validateServiceKey(ctx context.Context, expected string) bool {
	if expected == "" {
		return false
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}
	keys := md.Get(serviceKeyHeader)
	return len(keys) > 0 && keys[0] == expected
}

func authenticateUser(ctx context.Context, userClient userpb.UserServiceClient) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata ausente")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", status.Error(codes.Unauthenticated, "token ausente")
	}

	token := strings.TrimPrefix(tokens[0], "Bearer ")

	user, err := userClient.ValidateToken(ctx, &userpb.ValidateTokenRequest{Token: token})
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "token inválido")
	}

	return user.UserId, nil
}
