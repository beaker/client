package client

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/beaker/client/api"
)

// WhoAmI returns a client's active user.
func (c *Client) WhoAmI(ctx context.Context) (*api.UserDetail, error) {
	uri := path.Join("/api/v3/user")
	resp, err := c.sendRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	if err := errorFromResponse(resp); err != nil {
		return nil, err
	}

	var user api.UserDetail
	if err := parseResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GenerateToken creates a new token for authentication.
func (c *Client) GenerateToken(ctx context.Context) (string, error) {
	resp, err := c.sendRequest(ctx, http.MethodPost, "/api/v3/auth/tokens", nil, nil)
	if err != nil {
		return "", err
	}
	defer safeClose(resp.Body)

	if err := errorFromResponse(resp); err != nil {
		return "", err
	}

	var token string
	for _, cookie := range resp.Cookies() {
		if cookie.Name != "User-Token" {
			continue
		}
		token = cookie.Value
	}
	if token == "" {
		return "", fmt.Errorf("token not found in response")
	}
	return token, nil
}
