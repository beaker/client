package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// UserHandle provides operations on a user.
type UserHandle struct {
	client *Client
	id     string
}

// User gets a handle for a user by name or ID. The returned handle is
// guaranteed throughout its lifetime to refer to the same object, even if that
// object is later renamed.
func (c *Client) User(ctx context.Context, reference string) (*UserHandle, error) {
	// user is a top level resource
	if err := validateRef(reference, 1); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/users", reference)
	resp, err := c.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var user api.UserDetail
	if err := parseResponse(resp, &user); err != nil {
		return nil, err
	}

	return &UserHandle{client: c, id: user.ID}, nil
}

// ID returns a user's stable, unique ID.
func (h *UserHandle) ID() string {
	return h.id
}

// Get retrieves a user's details.
func (h *UserHandle) Get(ctx context.Context) (*api.UserDetail, error) {
	uri := path.Join("/api/v3/users", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.UserDetail
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &body, nil
}
