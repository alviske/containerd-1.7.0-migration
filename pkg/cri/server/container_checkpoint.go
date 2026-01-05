/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package server

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/runtime/v2/runc/options"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// TODO(KE): checkpoint implementation.
func (c *criService) CheckpointContainer(ctx context.Context, r *runtime.CheckpointContainerRequest) (res *runtime.CheckpointContainerResponse, err error) {
	container, err := c.containerStore.Get(r.GetContainerId())
	if err != nil {
		return nil, fmt.Errorf("an error occurred when try to find container %q: %w", r.GetContainerId(), err)
	}

	state := container.Status.Get().State()
	if state != runtime.ContainerState_CONTAINER_RUNNING {
		return nil, fmt.Errorf("container %q is in %s state. Only %s containers can be checkpointed",
			r.GetContainerId(), criContainerStateToString(state), criContainerStateToString(runtime.ContainerState_CONTAINER_RUNNING))
	}
	task, err := container.Container.Task(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get task for container %q: %w", r.GetContainerId(), err)
	}

	opts := []containerd.CheckpointTaskOpts{withCheckpointOpts(r.GetLocation(), r.GetParentPath())}
	log.G(ctx).Infof("Stop(Checkpoint) Time: %v.", time.Now().UnixMilli())
	_, err = task.Checkpoint(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("checkpoint container %q failed: %w", r.GetContainerId(), err)
	}
	return &runtime.CheckpointContainerResponse{}, nil
}

func withCheckpointOpts(location, parentPath string) containerd.CheckpointTaskOpts {
	return func(r *containerd.CheckpointTaskInfo) error {
		r.Name = location
		exit := true
		if r.Options == nil {
			r.Options = &options.CheckpointOptions{}
		}
		opts, _ := r.Options.(*options.CheckpointOptions)
		opts.WorkPath = location
		opts.ImagePath = location
		opts.ParentPath = parentPath
		opts.Exit = exit
		return nil
	}
}
