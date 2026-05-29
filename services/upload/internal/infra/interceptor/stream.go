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

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthStreamInterceptor(userClient userpb.UserServiceClient) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("StreamInterceptor chamado")

		md, ok := metadata.FromIncomingContext(stream.Context())
		log.Printf("metadata ok: %v", ok)
		log.Printf("metadata: %v", md)
		if !ok {
			return status.Error(codes.Unauthenticated, "metadata ausente")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return status.Error(codes.Unauthenticated, "token ausente")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")
		log.Printf("tokens: %v", tokens)

		user, err := userClient.ValidateToken(context.Background(), &userpb.ValidateTokenRequest{Token: token})
		if err != nil {
			log.Printf("ValidateToken err: %v", err)
			return status.Error(codes.Unauthenticated, "token inválido")
		}

		ctx := context.WithValue(stream.Context(), UserIDKey, user.UserId)
		return handler(srv, &wrappedStream{stream, ctx})
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
