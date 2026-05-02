package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/paknahad/docker-log/internal/domain"
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

func TestClientOpenContainerLogs(t *testing.T) {
	reader := io.NopCloser(strings.NewReader("ready\n"))
	api := &fakeContainerAPI{
		inspect:   dockertypes.ContainerJSON{Config: &dockercontainer.Config{Tty: true}},
		logReader: reader,
	}
	client := NewClientWithAPI(api)

	got, err := client.OpenContainerLogs(context.Background(), domain.Container{ID: "abc123", Name: "api"})
	if err != nil {
		t.Fatalf("OpenContainerLogs() error = %v", err)
	}
	body, err := io.ReadAll(got)
	if err != nil {
		t.Fatalf("read logs: %v", err)
	}
	if string(body) != "ready\n" {
		t.Fatalf("logs = %q, want ready newline", body)
	}
	if api.logContainerID != "abc123" {
		t.Fatalf("ContainerLogs() container = %q, want abc123", api.logContainerID)
	}
	if api.inspectContainerID != "abc123" {
		t.Fatalf("ContainerInspect() container = %q, want abc123", api.inspectContainerID)
	}
	if !api.logOptions.Follow {
		t.Fatal("ContainerLogs() Follow = false, want true")
	}
	if !api.logOptions.ShowStdout {
		t.Fatal("ContainerLogs() ShowStdout = false, want true")
	}
	if !api.logOptions.ShowStderr {
		t.Fatal("ContainerLogs() ShowStderr = false, want true")
	}
	if api.logOptions.Tail != "0" {
		t.Fatalf("ContainerLogs() Tail = %q, want 0 for live-only streams", api.logOptions.Tail)
	}
}

func TestClientOpenContainerLogsDemultiplexesDockerFrames(t *testing.T) {
	var framed bytes.Buffer
	stdout := stdcopy.NewStdWriter(&framed, stdcopy.Stdout)
	stderr := stdcopy.NewStdWriter(&framed, stdcopy.Stderr)
	if _, err := stdout.Write([]byte("out one\n")); err != nil {
		t.Fatalf("write stdout frame: %v", err)
	}
	if _, err := stderr.Write([]byte("err one\n")); err != nil {
		t.Fatalf("write stderr frame: %v", err)
	}

	api := &fakeContainerAPI{
		inspect:   dockertypes.ContainerJSON{Config: &dockercontainer.Config{Tty: false}},
		logReader: io.NopCloser(bytes.NewReader(framed.Bytes())),
	}
	client := NewClientWithAPI(api)

	got, err := client.OpenContainerLogs(context.Background(), domain.Container{ID: "abc123", Name: "api"})
	if err != nil {
		t.Fatalf("OpenContainerLogs() error = %v", err)
	}
	body, err := io.ReadAll(got)
	if err != nil {
		t.Fatalf("read logs: %v", err)
	}

	if string(body) != "out one\nerr one\n" {
		t.Fatalf("logs = %q, want demultiplexed stdout/stderr payloads", body)
	}
}

func TestClientOpenContainerLogsWrapsErrors(t *testing.T) {
	api := &fakeContainerAPI{logErr: errors.New("permission denied")}
	client := NewClientWithAPI(api)

	_, err := client.OpenContainerLogs(context.Background(), domain.Container{ID: "abc123", Name: "api"})
	if err == nil {
		t.Fatal("OpenContainerLogs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "open Docker logs for api") {
		t.Fatalf("error = %q, want log stream context", err)
	}
	if !errors.Is(err, api.logErr) {
		t.Fatalf("error does not wrap original error %v", api.logErr)
	}
}

func TestClientOpenContainerLogsWrapsInspectErrors(t *testing.T) {
	api := &fakeContainerAPI{inspectErr: errors.New("not found")}
	client := NewClientWithAPI(api)

	_, err := client.OpenContainerLogs(context.Background(), domain.Container{ID: "abc123", Name: "api"})
	if err == nil {
		t.Fatal("OpenContainerLogs() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "inspect Docker container api") {
		t.Fatalf("error = %q, want inspect context", err)
	}
	if !errors.Is(err, api.inspectErr) {
		t.Fatalf("error does not wrap original error %v", api.inspectErr)
	}
	if api.logContainerID != "" {
		t.Fatalf("ContainerLogs() container = %q, want no log call after inspect failure", api.logContainerID)
	}
}

type fakeContainerAPI struct {
	containers         []dockertypes.Container
	err                error
	options            dockercontainer.ListOptions
	inspectContainerID string
	inspect            dockertypes.ContainerJSON
	inspectErr         error
	logContainerID     string
	logOptions         dockercontainer.LogsOptions
	logReader          io.ReadCloser
	logErr             error
}

func (f *fakeContainerAPI) ContainerList(ctx context.Context, options dockercontainer.ListOptions) ([]dockertypes.Container, error) {
	f.options = options
	if f.err != nil {
		return nil, f.err
	}
	return f.containers, nil
}

func (f *fakeContainerAPI) ContainerInspect(ctx context.Context, containerID string) (dockertypes.ContainerJSON, error) {
	f.inspectContainerID = containerID
	if f.inspectErr != nil {
		return dockertypes.ContainerJSON{}, f.inspectErr
	}
	return f.inspect, nil
}

func (f *fakeContainerAPI) ContainerLogs(ctx context.Context, containerID string, options dockercontainer.LogsOptions) (io.ReadCloser, error) {
	f.logContainerID = containerID
	f.logOptions = options
	if f.logErr != nil {
		return nil, f.logErr
	}
	return f.logReader, nil
}
