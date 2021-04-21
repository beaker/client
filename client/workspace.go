package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/beaker/client/api"
)

type WorkspaceHandle struct {
	client *Client
	id     string
	// TODO: Support name-based references too when API support is available
}

func (c *Client) CreateWorkspace(
	ctx context.Context,
	spec api.WorkspaceSpec,
) (*WorkspaceHandle, error) {
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/workspaces", nil, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var workspace api.Workspace
	if err := parseResponse(resp, &workspace); err != nil {
		return nil, err
	}

	return &WorkspaceHandle{client: c, id: workspace.ID}, nil
}

type ListWorkspaceOptions struct {
	Archived *bool
	Cursor   string
	Text     string
}

func (c *Client) ListWorkspaces(
	ctx context.Context,
	org string,
	opts *ListWorkspaceOptions,
) ([]api.Workspace, string, error) {
	if opts == nil {
		opts = &ListWorkspaceOptions{}
	}

	query := url.Values{}
	query.Add("org", org)
	query.Add("cursor", opts.Cursor)
	if opts.Archived != nil {
		query.Add("archived", strconv.FormatBool(*opts.Archived))
	}
	if opts.Text != "" {
		query.Add("q", opts.Text)
	}

	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, "/api/v3/workspaces", query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.WorkspacePage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}

	return result.Data, result.NextCursor, nil
}

// Workspace gets a handle for a workspace by name or ID. The returned handle is
// guaranteed throughout its lifetime to refer to the same object, even if that
// object is later renamed.
func (c *Client) Workspace(ctx context.Context, reference string) (*WorkspaceHandle, error) {
	workspace, err := getWorkspace(ctx, c, reference)
	if err != nil {
		return nil, err
	}
	return &WorkspaceHandle{client: c, id: workspace.ID}, nil
}

// ID returns a workspace's stable, unique ID.
func (h *WorkspaceHandle) ID() string {
	return h.id
}

// Get retrieves a task's details.
func (h *WorkspaceHandle) Get(ctx context.Context) (*api.Workspace, error) {
	return getWorkspace(ctx, h.client, h.id)
}

func getWorkspace(ctx context.Context, c *Client, reference string) (*api.Workspace, error) {
	if strings.Count(reference, "/") != 1 {
		return nil, fmt.Errorf("%q isn't a valid workspace name, expected \"account/workspace\"", reference)
	}

	path := path.Join("/api/v3/workspaces", reference)
	resp, err := c.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var workspace api.Workspace
	if err := parseResponse(resp, &workspace); err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (h *WorkspaceHandle) Transfer(ctx context.Context, ids ...string) error {
	body := api.WorkspaceTransferSpec{IDs: ids}
	path := path.Join("/api/v3/workspaces", h.id, "transfer")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

type ListDatasetOptions struct {
	Cursor        string
	ResultsOnly   *bool
	CommittedOnly *bool
	Text          string
}

func (h *WorkspaceHandle) Datasets(
	ctx context.Context,
	opts *ListDatasetOptions,
) ([]api.Dataset, string, error) {
	if opts == nil {
		opts = &ListDatasetOptions{}
	}

	query := url.Values{}
	query.Add("cursor", opts.Cursor)
	if opts.ResultsOnly != nil {
		query.Add("results", strconv.FormatBool(*opts.ResultsOnly))
	}
	if opts.CommittedOnly != nil {
		query.Add("committed", strconv.FormatBool(*opts.CommittedOnly))
	}
	if opts.Text != "" {
		query.Add("q", opts.Text)
	}

	path := path.Join("/api/v3/workspaces", h.id, "datasets")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.DatasetPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}

	return result.Data, result.NextCursor, nil
}

// ExperimentOpts allows a caller to set options when creating an experiment.
// All fields are optional and default to reasonable values.
type ExperimentOpts struct {
	// Name is a human-friendly identifier for display.
	Name string

	// AuthorToken may be set to an API token to attribute an experiment. If
	// omitted, the author defaults to the bearer token set on the client.
	AuthorToken string
}

// CreateExperiment creates a new experiment in the calling workspace.
func (h *WorkspaceHandle) CreateExperiment(
	ctx context.Context,
	spec *api.ExperimentSpecV2,
	opts *ExperimentOpts,
) (*api.Experiment, error) {
	var query url.Values
	if opts != nil && opts.Name != "" {
		query = url.Values{"name": {opts.Name}}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(spec); err != nil {
		return nil, err
	}

	path := path.Join("/api/v3/workspaces", h.id, "experiments")
	req, err := h.client.newRetryableRequest(http.MethodPost, path, query, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if opts != nil && opts.AuthorToken != "" {
		req.Header.Set(api.HeaderAuthor, opts.AuthorToken)
	}

	resp, err := newRetryableClient(&http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}, h.client.HTTPResponseHook).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Experiment
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type ListExperimentOptions struct {
	Cursor string
	Text   string
}

func (h *WorkspaceHandle) Experiments(
	ctx context.Context,
	opts *ListExperimentOptions,
) ([]api.Experiment, string, error) {
	if opts == nil {
		opts = &ListExperimentOptions{}
	}

	query := url.Values{}
	query.Add("cursor", opts.Cursor)
	if opts.Text != "" {
		query.Add("q", opts.Text)
	}

	path := path.Join("/api/v3/workspaces", h.id, "experiments")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.ExperimentPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}

	return result.Data, result.NextCursor, nil
}

type ListGroupOptions struct {
	Cursor string
	Text   string
}

func (h *WorkspaceHandle) Groups(
	ctx context.Context,
	opts *ListGroupOptions,
) ([]api.Group, string, error) {
	if opts == nil {
		opts = &ListGroupOptions{}
	}

	query := url.Values{}
	query.Add("cursor", opts.Cursor)
	if opts.Text != "" {
		query.Add("q", opts.Text)
	}

	path := path.Join("/api/v3/workspaces", h.id, "groups")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.GroupPage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}

	return result.Data, result.NextCursor, nil
}

type ListImageOptions struct {
	Cursor string
	Text   string
}

func (h *WorkspaceHandle) Images(
	ctx context.Context,
	opts *ListImageOptions,
) ([]api.Image, string, error) {
	if opts == nil {
		opts = &ListImageOptions{}
	}

	query := url.Values{}
	query.Add("cursor", opts.Cursor)
	if opts.Text != "" {
		query.Add("q", opts.Text)
	}

	path := path.Join("/api/v3/workspaces", h.id, "images")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	defer safeClose(resp.Body)

	var result api.ImagePage
	if err := parseResponse(resp, &result); err != nil {
		return nil, "", err
	}

	return result.Data, result.NextCursor, nil
}

// SetName sets a workspace's name.
// References to the workspace's contents using the old name will stop working.
func (h *WorkspaceHandle) SetName(ctx context.Context, name string) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Name: &name})
}

// SetDescription sets a workspace's description.
func (h *WorkspaceHandle) SetDescription(ctx context.Context, desc string) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Description: &desc})
}

// SetArchived sets the archival status of a workspace.
// Archived workspaces are read-only.
func (h *WorkspaceHandle) SetArchived(ctx context.Context, archive bool) error {
	return h.patchWorkspace(ctx, api.WorkspacePatchSpec{Archive: &archive})
}

func (h *WorkspaceHandle) patchWorkspace(ctx context.Context, spec api.WorkspacePatchSpec) error {
	path := path.Join("/api/v3/workspaces", h.id)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, spec)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (h *WorkspaceHandle) Permissions(ctx context.Context) (*api.WorkspacePermissionSummary, error) {
	path := path.Join("/api/v3/workspaces", h.id, "auth")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.WorkspacePermissionSummary
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *WorkspaceHandle) SetPermissions(ctx context.Context, patch api.WorkspacePermissionPatch) error {
	path := path.Join("/api/v3/workspaces", h.id, "auth")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, patch)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (h *WorkspaceHandle) ListSecrets(ctx context.Context) ([]api.Secret, error) {
	path := path.Join("/api/v3/workspaces", h.id, "secrets")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Secrets
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (h *WorkspaceHandle) GetSecret(ctx context.Context, name string) (*api.Secret, error) {
	path := path.Join("/api/v3/workspaces", h.id, "secrets", name)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Secret
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *WorkspaceHandle) DeleteSecret(ctx context.Context, name string) error {
	path := path.Join("/api/v3/workspaces", h.id, "secrets", name)
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (h *WorkspaceHandle) PutSecret(
	ctx context.Context,
	name string,
	value []byte,
) (*api.Secret, error) {
	path := path.Join("/api/v3/workspaces", h.id, "secrets", name, "value")
	req, err := h.client.newRetryableRequest(http.MethodPut, path, nil, bytes.NewReader(value))
	if err != nil {
		return nil, err
	}

	resp, err := newRetryableClient(&http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}, h.client.HTTPResponseHook).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Secret
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *WorkspaceHandle) ReadSecret(ctx context.Context, name string) ([]byte, error) {
	path := path.Join("/api/v3/workspaces", h.id, "secrets", name, "value")
	req, err := h.client.newRetryableRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := newRetryableClient(&http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}, h.client.HTTPResponseHook).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	if err := errorFromResponse(resp); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
