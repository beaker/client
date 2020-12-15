package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// NodeHandle provides access to a single node.
type NodeHandle struct {
	client *Client
	id     string
}

// Node gets a handle for an node by ID. The node is not resolved
// and not guaranteed to exist.
func (c *Client) Node(id string) *NodeHandle {
	return &NodeHandle{client: c, id: id}
}

// ListExecutions retrieves all executions that are assigned to the node.
func (h *NodeHandle) ListExecutions(ctx context.Context) (*api.Executions, error) {
	path := path.Join("/api/v3/nodes", h.id, "executions")
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Executions
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AssignExecutions lists all executions on a node and assigns a new one if it has none.
func (h *NodeHandle) AssignExecutions(ctx context.Context, resources *api.NodeResources) (*api.Executions, error) {
	path := path.Join("/api/v3/nodes", h.id, "executions")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, resources)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Executions
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes the node (marks it terminated)
func (h *NodeHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/nodes", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Patch updates the fields of a node.
func (h *NodeHandle) Patch(ctx context.Context, patch *api.NodePatchSpec) error {
	path := path.Join("/api/v3/nodes", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, patch)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
