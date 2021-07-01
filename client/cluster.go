package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/beaker/client/api"
)

// CreateCluster creates a new dynamic compute cluster on which to run experiments.
func (c *Client) CreateCluster(
	ctx context.Context,
	account string,
	spec api.ClusterSpec,
) (*api.Cluster, error) {
	path := path.Join("/api/v3/clusters", url.PathEscape(account))
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, path, nil, spec)
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

type ListClusterOptions struct {
	Limit  int64
	Cursor string
}

// ListClusters enumerates all clusters owned by an account. If provided, only clusters
// that match the constraints in opts are returned.
func (c *Client) ListClusters(
	ctx context.Context,
	account string,
	opts *ListClusterOptions,
) ([]api.Cluster, string, error) {
	if opts == nil {
		opts = &ListClusterOptions{}
	}

	path := path.Join("/api/v3/clusters", url.PathEscape(account))
	query := url.Values{}
	query.Add("cursor", opts.Cursor)
	if opts.Limit > 0 {
		query.Add("limit", strconv.FormatInt(opts.Limit, 10))
	}

	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.ClusterPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}
	return result.Data, result.NextCursor, nil
}

// Cluster gets a handle for a cluster by name or ID. The reference is not resolved.
func (c *Client) Cluster(reference string) *ClusterHandle {
	return &ClusterHandle{client: c, ref: reference}
}

// ClusterHandle provides access to a single cluster.
type ClusterHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *ClusterHandle) Ref() string {
	return h.ref
}

// Get retrieves a clusters details.
func (h *ClusterHandle) Get(ctx context.Context) (*api.Cluster, error) {
	if err := validateClusterRef(h.ref); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", h.ref)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
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

// Patch updates a cluster's details.
func (h *ClusterHandle) Patch(ctx context.Context, patch *api.ClusterPatch) (*api.Cluster, error) {
	if err := validateClusterRef(h.ref); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", h.ref)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, patch)
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

// Terminate invalidates a cluster and frees its name for reuse.
//
// New tasks cannot be created on the cluster, but existing scheduled tasks will
// be allowed to complete.
func (h *ClusterHandle) Terminate(ctx context.Context) error {
	if err := validateClusterRef(h.ref); err != nil {
		return err
	}

	path := path.Join("/api/v3/clusters", h.ref)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// CreateNode is meant for internal use only.
func (h *ClusterHandle) CreateNode(ctx context.Context, spec api.NodeSpec) (*api.Node, error) {
	if err := validateClusterRef(h.ref); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", h.ref, "nodes")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPost, path, nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)
	var result api.Node
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListClusterNodes enumerates all active nodes within a cluster.
// TODO: Make this return an iterator.
func (h *ClusterHandle) ListClusterNodes(ctx context.Context) ([]api.Node, error) {
	if err := validateClusterRef(h.ref); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/clusters", h.ref, "nodes")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.NodePage
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

type ExecutionFilters struct {
	Scheduled *bool
}

// ListExecutions enumerates all active or pending tasks on a cluster.
// TODO: Make this return an iterator.
func (h *ClusterHandle) ListExecutions(
	ctx context.Context,
	opts *ExecutionFilters,
) ([]api.Execution, error) {
	if err := validateClusterRef(h.ref); err != nil {
		return nil, err
	}

	v := url.Values{}
	if opts != nil {
		if opts.Scheduled != nil {
			v.Add("scheduled", strconv.FormatBool(*opts.Scheduled))
		}
	}

	path := path.Join("/api/v3/clusters", h.ref, "executions")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Executions
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// PatchExecution updates an execution within the cluster.
func (h *ClusterHandle) PatchExecution(
	ctx context.Context,
	execution string,
	spec api.ExecutionPatchSpec,
) error {
	if err := validateClusterRef(h.ref); err != nil {
		return err
	}

	path := path.Join("/api/v3/clusters", h.ref, "executions", execution)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, spec)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Validate a cluster reference appears to be of the correct shape. This helps
// catch errors relating to missing delimiters so we can show consistent errors.
func validateClusterRef(ref string) error {
	if strings.Count(ref, "/") != 1 {
		return fmt.Errorf("%q isn't a valid cluster name, expected \"account/cluster\"", ref)
	}
	return nil
}
