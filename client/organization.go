package client

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"

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

// VerifyOrgExists validates existence of an organization.
func (c *Client) VerifyOrgExists(ctx context.Context, org string) error {
	resp, err := c.sendRequest(ctx, http.MethodGet, path.Join("/api/v3/orgs", org), nil, nil)
	if err != nil {
		return errors.WithMessage(err, "could not resolve organization "+org)
	}
	defer safeClose(resp.Body)

	return errors.WithMessage(errorFromResponse(resp), org)
}
