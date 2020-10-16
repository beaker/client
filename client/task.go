package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// TaskHandle provides operations on a task.
type TaskHandle struct {
	client *Client
	id     string
}

// Task gets a handle for a task by name or ID. The returned handle is
// guaranteed throughout its lifetime to refer to the same object, even if that
// object is later renamed, however, the task is not resolved and may infact
// not exist/already be deleted.
func (c *Client) Task(reference string) *TaskHandle {
	return &TaskHandle{client: c, id: reference}
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
