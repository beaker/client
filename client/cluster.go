package client

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/beaker/client/api"
)

// ClusterHandle provides access to a single cluster.
type ClusterHandle struct {
	client *Client
	name   string
}

// CreateCluster creates a new dynamic compute cluster on which to run experiments.
func (c *Client) CreateCluster(
	ctx context.Context,
	account string,
	spec api.ClusterSpec,
) (*api.Cluster, error) {
	if err := validateRef(account, 1); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", account)
	resp, err := c.sendRequest(ctx, http.MethodPost, path, nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Cluster
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListClusters enumerates all clusters within a galaxy.
// TODO: Make this return an iterator.
// TODO: Include galaxy, expiration filter.
func (c *Client) ListClusters(
	ctx context.Context,
	account string,
) ([]api.Cluster, error) {
	if err := validateRef(account, 1); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", account)
	resp, err := c.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.ClusterPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// Cluster gets a handle for a cluster by name or ID. The cluster is not resolved
// and not guaranteed to exist.
func (c *Client) Cluster(name string) *ClusterHandle {
	return &ClusterHandle{client: c, name: name}
}

// Name returns the name or ID used in creation of the cluster handle.
func (h *ClusterHandle) Name() string {
	return h.name
}

// Get retrieves a clusters details.
func (h *ClusterHandle) Get(ctx context.Context) (*api.Cluster, error) {
	if err := validateRef(h.name, 2); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", h.name)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Cluster
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Extend resets a cluster's time-to-live (TTL), returning the new expiration time.
func (h *ClusterHandle) Extend(ctx context.Context) (time.Time, error) {
	if err := validateRef(h.name, 2); err != nil {
		return time.Time{}, err
	}

	path := path.Join("/api/v3/clusters", h.name, "extend")
	resp, err := h.client.sendRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return time.Time{}, err
	}
	defer safeClose(resp.Body)

	var result api.Cluster
	if err := parseResponse(resp, &result); err != nil {
		return time.Time{}, err
	}
	if result.Expiration == nil {
		return time.Time{}, nil
	}
	return *result.Expiration, nil
}

// Delete terminates a cluster.
func (h *ClusterHandle) Delete(ctx context.Context) error {
	if err := validateRef(h.name, 2); err != nil {
		return err
	}

	path := path.Join("/api/v3/clusters", h.name)
	resp, err := h.client.sendRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
