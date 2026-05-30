package upload

import (
	"context"

	uploadpb "transcode-service/pkg/uploadpb"
)

// GRPCClient adapta o cliente gRPC do upload para o use case.
type GRPCClient struct {
	client uploadpb.UploadServiceClient
}

func NewGRPCClient(client uploadpb.UploadServiceClient) *GRPCClient {
	return &GRPCClient{client: client}
}

func (c *GRPCClient) UpdateStatus(ctx context.Context, videoID string, status uploadpb.VideoStatus) error {
	_, err := c.client.UpdateStatus(ctx, &uploadpb.UpdateStatusRequest{
		Id:     videoID,
		Status: status,
	})
	return err
}
