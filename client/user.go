package client

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/beaker/client/api"
)

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

// User gets a handle for a user by name or ID. The reference is not resolved.
func (c *Client) User(reference string) *UserHandle {
	return &UserHandle{client: c, ref: reference}
}

// UserHandle provides operations on a user.
type UserHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *UserHandle) Ref() string {
	return h.ref
}

// Get retrieves a user's details.
func (h *UserHandle) Get(ctx context.Context) (*api.UserDetail, error) {
	uri := path.Join("/api/v3/users", url.PathEscape(h.ref))
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
