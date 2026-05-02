package docker

import (
	"context"
	"errors"
	"strings"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
)

func TestClientListRunningContainers(t *testing.T) {
	api := &fakeContainerAPI{
		containers: []dockertypes.Container{
			{
				ID:     "abc123",
				Names:  []string{"/api", "/compose-api-1"},
				Image:  "example/api:latest",
				Status: "Up 2 minutes",
			},
			{
				ID:     "def456",
				Names:  []string{"worker"},
				Image:  "example/worker:latest",
				Status: "Up 1 minute",
			},
		},
	}
	client := NewClientWithAPI(api)

	containers, err := client.ListRunningContainers(context.Background())
	if err != nil {
		t.Fatalf("ListRunningContainers() error = %v", err)
	}

	if len(containers) != 2 {
		t.Fatalf("len(containers) = %d, want 2", len(containers))
	}
	if containers[0].ID != "abc123" {
		t.Fatalf("containers[0].ID = %q, want abc123", containers[0].ID)
	}
	if containers[0].Name != "api" {
		t.Fatalf("containers[0].Name = %q, want api", containers[0].Name)
	}
	if containers[0].Image != "example/api:latest" {
		t.Fatalf("containers[0].Image = %q, want example/api:latest", containers[0].Image)
	}
	if containers[0].Status != "Up 2 minutes" {
		t.Fatalf("containers[0].Status = %q, want Up 2 minutes", containers[0].Status)
	}
	if containers[1].Name != "worker" {
		t.Fatalf("containers[1].Name = %q, want worker", containers[1].Name)
	}
	if api.options.All {
		t.Fatalf("ContainerList() All = true, want false for running containers only")
	}
}

func TestClientListRunningContainersWrapsErrors(t *testing.T) {
	api := &fakeContainerAPI{err: errors.New("daemon unavailable")}
	client := NewClientWithAPI(api)

	_, err := client.ListRunningContainers(context.Background())
	if err == nil {
		t.Fatal("ListRunningContainers() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "list running Docker containers") {
		t.Fatalf("error = %q, want discovery context", err)
	}
	if !errors.Is(err, api.err) {
		t.Fatalf("error does not wrap original error %v", api.err)
	}
}

type fakeContainerAPI struct {
	containers []dockertypes.Container
	err        error
	options    dockercontainer.ListOptions
}

func (f *fakeContainerAPI) ContainerList(ctx context.Context, options dockercontainer.ListOptions) ([]dockertypes.Container, error) {
	f.options = options
	if f.err != nil {
		return nil, f.err
	}
	return f.containers, nil
}
