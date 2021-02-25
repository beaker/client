package client

import (
	"context"
	"net/http"
	"net/url"
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
	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
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
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, uri, nil, nil)
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

// ListUsers gets all users.
func (c *Client) ListUsers(
	ctx context.Context,
	cursor string,
) ([]api.UserDetail, string, error) {
	query := url.Values{}
	query.Add("cursor", cursor)
	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, "/api/v3/admin/users", query, nil)
	if err != nil {
		return nil, "", err
	}

	var result api.UserPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}
	return result.Data, result.NextCursor, nil
}
