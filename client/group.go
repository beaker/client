package client

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/beaker/client/api"
)

// CreateGroup creates a new group with an optional name.
func (c *Client) CreateGroup(ctx context.Context, spec api.GroupSpec) (*GroupHandle, error) {
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/groups", nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.CreateGroupResponse
	if err = parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &GroupHandle{client: c, ref: body.ID}, nil
}

// Group gets a handle for a group by name or ID. The reference is not resolved.
func (c *Client) Group(reference string) *GroupHandle {
	return &GroupHandle{client: c, ref: reference}
}

// GroupHandle provides operations on a group.
type GroupHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *GroupHandle) Ref() string {
	return h.ref
}

// Get retrieves a group's details.
func (h *GroupHandle) Get(ctx context.Context) (*api.Group, error) {
	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.Group
	if err = parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// SetName sets a group's name.
func (h *GroupHandle) SetName(ctx context.Context, name string) error {
	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	body := api.GroupPatchSpec{Name: &name}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// SetDescription sets a group's description.
func (h *GroupHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	body := api.GroupPatchSpec{Description: &description}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Experiments returns the IDs of all experiments within a group.
func (h *GroupHandle) Experiments(ctx context.Context) ([]string, error) {
	path := path.Join("/api/v3/groups", url.PathEscape(h.ref), "experiments")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body []string
	if err = parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return body, nil
}

// AddExperiments adds experiments by name or ID to a group.
func (h *GroupHandle) AddExperiments(ctx context.Context, experiments []string) error {
	if len(experiments) == 0 {
		return nil
	}

	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	body := api.GroupPatchSpec{AddExperiments: experiments}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)

	return errorFromResponse(resp)
}

// RemoveExperiments removes experiments by name or ID from a group.
func (h *GroupHandle) RemoveExperiments(ctx context.Context, experiments []string) error {
	if len(experiments) == 0 {
		return nil
	}

	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	body := api.GroupPatchSpec{RemoveExperiments: experiments}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)

	return errorFromResponse(resp)
}

// Delete removes a group and its contents.
func (h *GroupHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/groups", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)

	return errorFromResponse(resp)
}
