package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	retryable "github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/beaker/client/api"
)

// TaskHandle provides operations on a task.
type TaskHandle struct {
	client *Client
	id     string
}

// Task gets a handle for a task by name or ID. The returned handle is
// guaranteed throughout its lifetime to refer to the same object, even if that
// object is later renamed.
func (c *Client) Task(ctx context.Context, reference string) (*TaskHandle, error) {
	id, err := c.resolveRef(ctx, "/api/v3/tasks", reference)
	if err != nil {
		return nil, errors.WithMessage(err, "could not resolve task reference "+reference)
	}

	return &TaskHandle{client: c, id: id}, nil
}

// ID returns a task's stable, unique ID.
func (h *TaskHandle) ID() string {
	return h.id
}

// Get retrieves a task's details.
func (h *TaskHandle) Get(ctx context.Context) (*api.Task, error) {
	path := path.Join("/api/v3/tasks", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var task api.Task
	if err := parseResponse(resp, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// GetResults retrieves a task's results.
func (h *TaskHandle) GetResults(ctx context.Context) (*api.TaskResults, error) {
	path := path.Join("/api/v3/tasks", h.id, "results")
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var results api.TaskResults
	if err := parseResponse(resp, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// SetDescription sets a task's description.
func (h *TaskHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/tasks", h.id)
	body := api.TaskPatchSpec{Description: &description}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Stop cancels a task. If the task has already completed, this succeeds with no effect.
func (h *TaskHandle) Stop(ctx context.Context) error {
	path := path.Join("/api/v3/tasks", h.id)
	body := api.TaskPatchSpec{Cancel: true}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// GetLogs gets all logs for a task. Logs are in the form:
// {RFC3339 nano timestamp} {message}\n
func (h *TaskHandle) GetLogs(ctx context.Context) (io.ReadCloser, error) {
	path := path.Join("/api/v3/tasks", h.id, "logs")
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := errorFromResponse(resp); err != nil {
		safeClose(resp.Body)
		return nil, err
	}
	return resp.Body, nil
}

// PutLogs uploads a log chunk. Since is the time of the first log message in the chunk.
func (h *TaskHandle) PutLogs(ctx context.Context, logs io.Reader, since time.Time) error {
	// 0 represents task run. It's safe to ignore this since we can't have more than one run per task.
	path := fmt.Sprintf("/api/tasks/%s/logs/0/task-%s.log/upload", h.id, since.Format(time.RFC3339Nano))
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)

	var uploadLink api.TaskLogUploadLink
	if err := parseResponse(resp, &uploadLink); err != nil {
		return err
	}

	req, err := retryable.NewRequest(http.MethodPut, uploadLink.URL, logs)
	if err != nil {
		return err
	}

	resp, err = newRetryableClient(&http.Client{
		Timeout: 30 * time.Second,
	}, nil).Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.Errorf("log upload failed with status %d", resp.StatusCode)
	}
	return nil
}

// ExecutionHandle provides access to a single execution.
type ExecutionHandle struct {
	client *Client
	taskID string
	id     string
}

// Execution gets a handle for an execution by ID. The execution is not resolved
// and not guaranteed to exist.
func (h *TaskHandle) Execution(id string) *ExecutionHandle {
	return &ExecutionHandle{client: h.client, taskID: h.id, id: id}
}

// Get retrieves an execution's details.
func (h *ExecutionHandle) Get(ctx context.Context) (*api.Execution, error) {
	path := path.Join("/api/v3/tasks", h.taskID, "executions", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Execution
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PostStatus updates an execution's current status.
func (h *ExecutionHandle) PostStatus(ctx context.Context, status api.ExecStatusUpdate) error {
	path := path.Join("/api/v3/tasks", h.taskID, "executions", h.id, "status")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, status)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
