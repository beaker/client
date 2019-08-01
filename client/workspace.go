package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

type WorkspaceHandle struct {
	client *Client
	id     string
	// TODO: Support name-based references too when API support is available
}

func (c *Client) CreateWorkspace(ctx context.Context, spec api.WorkspaceSpec) (*WorkspaceHandle, error) {
	resp, err := c.sendRequest(ctx, http.MethodPost, "/api/v3/workspaces", nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var workspaceID api.CreateWorkspaceResponse
	if err := parseResponse(resp, &workspaceID); err != nil {
		return nil, err
	}

	return &WorkspaceHandle{client: c, id: workspaceID.ID}, nil
}

// Workspace gets a handle for a workspace by name or ID. The returned handle is
// guaranteed throughout its lifetime to refer to the same object, even if that
// object is later renamed.
func (c *Client) Workspace(ctx context.Context, reference string) (*WorkspaceHandle, error) {
	// TODO: c.resolveRef so reference can be either name or ID, once API endpoints accept names
	return &WorkspaceHandle{client: c, id: reference}, nil
}

// Get retrieves a task's details.
func (h *WorkspaceHandle) Get(ctx context.Context) (*api.Workspace, error) {
	path := path.Join("/api/v3/workspaces", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var workspace api.Workspace
	if err := parseResponse(resp, &workspace); err != nil {
		return nil, err
	}

	return &workspace, nil
}

// SetName sets a workspace's name.
// References to the workspace's contents using the old name will stop working.
func (h *WorkspaceHandle) SetName(ctx context.Context, name string) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Name: &name})
}

// SetDescription sets a workspace's description.
func (h *WorkspaceHandle) SetDescription(ctx context.Context, desc string) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Description: &desc})
}

// SetArchived sets the archival status of a workspace.
// Archived workspaces are read-only.
func (h *WorkspaceHandle) SetArchived(ctx context.Context, archive bool) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Archive: &archive})
}

func (h *WorkspaceHandle) patchWorkspace(ctx context.Context, spec api.WorkspacePatchSpec) error {
	path := path.Join("/api/v3/workspaces", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, spec)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
