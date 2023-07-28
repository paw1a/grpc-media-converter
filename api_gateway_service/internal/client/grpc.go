package client

import (
	"context"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"time"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewStorageServiceConn(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {
	//opts := []grpc_retry.CallOption{
	//	grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
	//	grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
	//	grpc_retry.WithMax(backoffRetries),
	//}

	storageServiceConn, err := grpc.DialContext(
		ctx,
		cfg.Grpc.StorageServicePort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}

	return storageServiceConn, nil
}
