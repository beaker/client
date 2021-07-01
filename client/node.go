package client

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/beaker/client/api"
)

// Node gets a handle for an node by ID. The id is not resolved.
func (c *Client) Node(id string) *NodeHandle {
	return &NodeHandle{client: c, id: id}
}

// NodeHandle provides access to a single node.
type NodeHandle struct {
	client *Client
	id     string
}

// Ref returns the name or ID with which a handle was created.
func (h *NodeHandle) Ref() string {
	return h.id
}

// Get information about a node.
func (h *NodeHandle) Get(ctx context.Context) (*api.Node, error) {
	path := path.Join("/api/v3/nodes", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var node api.Node
	if err := parseResponse(resp, &node); err != nil {
		return nil, err
	}

	return &node, nil
}

// ListExecutions retrieves all executions that are assigned to the node.
func (h *NodeHandle) ListExecutions(ctx context.Context) (*api.Executions, error) {
	path := path.Join("/api/v3/nodes", url.PathEscape(h.id), "executions")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
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
	path := path.Join("/api/v3/nodes", url.PathEscape(h.id), "executions")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPost, path, nil, resources)
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
	path := path.Join("/api/v3/nodes", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Patch updates the fields of a node.
func (h *NodeHandle) Patch(ctx context.Context, patch *api.NodePatchSpec) error {
	path := path.Join("/api/v3/nodes", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, patch)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
