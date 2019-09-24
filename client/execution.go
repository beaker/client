package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// ExecutionHandle provides access to a single execution.
type ExecutionHandle struct {
	client *Client
	id     string
}

// Execution gets a handle for an execution by ID. The execution is not resolved
// and not guaranteed to exist.
func (c *Client) Execution(id string) *ExecutionHandle {
	return &ExecutionHandle{client: c, id: id}
}

// Get retrieves an execution's details.
func (h *ExecutionHandle) Get(ctx context.Context) (*api.Execution, error) {
	path := path.Join("/api/v3/executions", h.id)
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
