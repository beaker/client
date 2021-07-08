package client

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/beaker/client/api"
)

// Task gets a handle for a task by ID. The id is not resolved.
func (c *Client) Task(id string) *TaskHandle {
	return &TaskHandle{client: c, id: id}
}

// TaskHandle provides operations on a task.
type TaskHandle struct {
	client *Client
	id     string
}

// Ref returns the name or ID with which a handle was created.
func (h *TaskHandle) Ref() string {
	return h.id
}

// Get retrieves a task's details.
func (h *TaskHandle) Get(ctx context.Context) (*api.Task, error) {
	path := path.Join("/api/v3/tasks", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
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
