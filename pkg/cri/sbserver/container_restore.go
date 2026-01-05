package sbserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (c *criService) RestoreContainer(ctx context.Context, r *runtime.RestoreContainerRequest) (res *runtime.RestoreContainerResponse, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckpointContainer not implemented")
}
