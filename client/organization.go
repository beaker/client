package client

import (
	"context"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// ListMyOrgs lists all orgs in which the caller is a member. The caller's
// account is inferred from the client's auth token.
func (c *Client) ListMyOrgs(ctx context.Context) ([]api.Organization, error) {
	resp, err := c.sendRequest(ctx, http.MethodGet, path.Join("/api/v3/user/orgs"), nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var orgs api.OrganizationPage
	if err := parseResponse(resp, &orgs); err != nil {
		return nil, err
	}

	// Cursor is not populated; assume the page contains all results.
	return orgs.Data, nil
}

// OrgHandle provides an interface for organization APIs.
type OrgHandle struct {
	client *Client
	ref    string
}

// Organization gets a handle for an organization by name or ID. The returned
// handle is not guaranteed to be valid.
func (c *Client) Organization(reference string) *OrgHandle {
	return &OrgHandle{client: c, ref: reference}
}

// Get retrieves an organization's details.
func (h *OrgHandle) Get(ctx context.Context) (*api.Organization, error) {
	path := path.Join("/api/v3/orgs", h.ref)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var org api.Organization
	if err := parseResponse(resp, &org); err != nil {
		return nil, err
	}

	return &org, nil
}

// ListMembers retrieves an organization's details.
// TODO: Make this return an iterator.
func (h *OrgHandle) ListMembers(
	ctx context.Context,
	cursor string,
) (users []api.UserDetail, next string, err error) {
	path := path.Join("/api/v3/orgs", h.ref, "members")
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var members api.UserPage
	if err := parseResponse(resp, &members); err != nil {
		return nil, "", err
	}

	return members.Data, members.NextCursor, nil
}

// GetMember returns details about a specific membership, if it exists.
func (h *OrgHandle) GetMember(ctx context.Context, account string) (*api.OrgMembership, error) {
	path := path.Join("/api/v3/orgs", h.ref, "members", account)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var member api.OrgMembership
	if err := parseResponse(resp, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// SetMember adds or updates the given account as a member of the org.
// Role must be "admin" or "member".
func (h *OrgHandle) SetMember(ctx context.Context, account string, role string) error {
	path := path.Join("/api/v3/orgs", h.ref, "members", account)
	resp, err := h.client.sendRequest(ctx, http.MethodPut, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// RemoveMember removes the given account from the org.
func (h *OrgHandle) RemoveMember(ctx context.Context, account string) error {
	path := path.Join("/api/v3/orgs", h.ref, "members", account)
	resp, err := h.client.sendRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
