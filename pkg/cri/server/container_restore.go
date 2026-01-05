package server

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/log"
	"github.com/pkg/errors"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (c *criService) RestoreContainer(ctx context.Context, r *runtime.RestoreContainerRequest) (res *runtime.RestoreContainerResponse, err error) {
	checkpointPath, _ := path.Split(strings.Split(r.GetLocation(), "/migration/")[1])
	if err := c.startContainer(ctx, r.GetContainerId(), checkpointPath, containerd.WithRestoreImagePath(r.GetLocation())); err != nil {
		return nil, errors.Wrap(err, "failed to restore container")
	}
	log.G(ctx).Infof("Restore Time: %v.", time.Now().UnixMilli())
	return &runtime.RestoreContainerResponse{}, nil
}
