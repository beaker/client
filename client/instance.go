package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// InstanceHandle provides access to a single instance.
type InstanceHandle struct {
	client *Client
	id     string
}

// Instance gets a handle for an instance by ID. The instance is not resolved
// and not guaranteed to exist.
func (c *Client) Instance(id string) *InstanceHandle {
	return &InstanceHandle{client: c, id: id}
}

// ListExecutions retrieves all executions that are assigned to the instance.
func (h *InstanceHandle) ListExecutions(ctx context.Context) (*api.InstanceExecutions, error) {
	path := path.Join("/api/v3/instances", h.id, "executions")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.InstanceExecutions
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
