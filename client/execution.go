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

// GetResults retrieves an execution's results.
func (h *ExecutionHandle) GetResults(ctx context.Context) (*api.ExecutionResults, error) {
	path := path.Join("/api/v3/executions", h.id, "results")
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var results api.ExecutionResults
	if err := parseResponse(resp, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// PostStatus updates an execution's current status.
func (h *ExecutionHandle) PostStatus(ctx context.Context, status api.ExecStatusUpdate) error {
	path := path.Join("/api/v3/executions", h.id, "status")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, status)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
