package interceptor

import (
	"context"

	userpb "upload-service/pkg/userpb"

	"google.golang.org/grpc"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// AuthStreamInterceptor autentica streams e adiciona o user_id ao contexto.
func AuthStreamInterceptor(userClient userpb.UserServiceClient) grpc.StreamServerInterceptor {
	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		userID, err := authenticateUser(stream.Context(), userClient)
		if err != nil {
			return err
		}

		ctx := context.WithValue(stream.Context(), UserIDKey, userID)
		return handler(srv, &wrappedStream{ServerStream: stream, ctx: ctx})
	}
}

// wrappedStream substitui o contexto original do stream.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context retorna o contexto enriquecido.
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
