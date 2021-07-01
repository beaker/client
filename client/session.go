package client

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/beaker/client/api"
)

// CreateSession creates an interactive Beaker session.
func (c *Client) CreateSession(ctx context.Context, spec api.SessionSpec) (*api.Session, error) {
	path := path.Join("/api/v3/sessions")
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, path, nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Session
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListSessionOpts specifies filters for listing sessions.
type ListSessionOpts struct {
	Node      *string
	Cluster   *string
	Finalized *bool
}

// ListSessions enumerates all sessions with optional filtering.
func (c *Client) ListSessions(
	ctx context.Context,
	opts *ListSessionOpts,
) ([]api.Session, error) {
	if opts == nil {
		opts = &ListSessionOpts{}
	}

	path := path.Join("/api/v3/sessions")
	query := url.Values{}
	if opts.Node != nil {
		query.Add("node", *opts.Node)
	}
	if opts.Cluster != nil {
		query.Add("cluster", *opts.Cluster)
	}
	if opts.Finalized != nil {
		query.Add("finalized", strconv.FormatBool(*opts.Finalized))
	}
	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result []api.Session
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Session gets a handle for a session by ID. The session is not resolved and
// not guaranteed to exist.
func (c *Client) Session(id string) *SessionHandle {
	return &SessionHandle{client: c, id: id}
}

// SessionHandle provides access to a single session.
type SessionHandle struct {
	client *Client
	id     string
}

// Ref returns the name or ID with which a handle was created.
func (h *SessionHandle) Ref() string {
	return h.id
}

// Get retrieves an session's details.
func (h *SessionHandle) Get(ctx context.Context) (*api.Session, error) {
	path := path.Join("/api/v3/sessions", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Session
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Patch updates a session.
func (h *SessionHandle) Patch(ctx context.Context, patch api.SessionPatch) (*api.Session, error) {
	path := path.Join("/api/v3/sessions", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, patch)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Session
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
