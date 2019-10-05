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
func (h *InstanceHandle) ListExecutions(ctx context.Context) (*api.ScheduledTasks, error) {
	path := path.Join("/api/v3/instances", h.id, "executions")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.ScheduledTasks
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetStatus retrieves all executions that are assigned to the instance.
func (h *InstanceHandle) SetStatus(ctx context.Context, status api.InstanceStatus) error {
	path := path.Join("/api/v3/instances", h.id, "status")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, &api.InstanceStatusSpec{
		Status: status,
	})
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Delete removes the instance (marks it terminated)
func (h *InstanceHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/instances", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
