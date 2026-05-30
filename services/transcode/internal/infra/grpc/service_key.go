package grpc

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const serviceKeyHeader = "x-service-key"

// ServiceKeyUnaryInterceptor envia a chave interna em chamadas ao upload-service.
func ServiceKeyUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		key := os.Getenv("INTERNAL_SERVICE_KEY")
		if key != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, serviceKeyHeader, key)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
