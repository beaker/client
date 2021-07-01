package client

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	fileheap "github.com/beaker/fileheap/client"

	"github.com/beaker/client/api"
)

// CreateDataset creates a new dataset with an optional name.
func (c *Client) CreateDataset(
	ctx context.Context,
	spec api.DatasetSpec,
	name string,
) (*DatasetHandle, error) {
	query := url.Values{}
	if name != "" {
		query.Set("name", name)
	}

	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/datasets", query, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.Dataset
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}

	return &DatasetHandle{client: c, ref: body.ID}, nil
}

// Dataset gets a handle for a dataset by name or ID. The reference is not resolved.
func (c *Client) Dataset(reference string) *DatasetHandle {
	return &DatasetHandle{client: c, ref: reference}
}

// DatasetHandle provides operations on a dataset.
type DatasetHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *DatasetHandle) Ref() string {
	return h.ref
}

// Get retrieves a dataset's details.
func (h *DatasetHandle) Get(ctx context.Context) (*api.Dataset, error) {
	uri := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.Dataset
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// Storage gets a client to access a dataset's backing storage. The returned
// client expires at the returned time and must be discarded and replaced.
func (h *DatasetHandle) Storage(ctx context.Context) (
	client *fileheap.DatasetRef,
	expiry time.Time,
	err error,
) {
	uri := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer safeClose(resp.Body)

	var body api.Dataset
	if err := parseResponse(resp, &body); err != nil {
		return nil, time.Time{}, err
	}

	fh, err := fileheap.New(body.Storage.Address, fileheap.WithToken(body.Storage.Token))
	if err != nil {
		return nil, time.Time{}, err
	}

	return fh.Dataset(body.Storage.ID), body.Storage.TokenExpires, nil
}

// SetName sets a dataset's name.
func (h *DatasetHandle) SetName(ctx context.Context, name string) error {
	path := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	body := api.DatasetPatchSpec{Name: &name}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// SetDescription sets a dataset's description.
func (h *DatasetHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	body := api.DatasetPatchSpec{Description: &description}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Commit finalizes a dataset, unblocking usage and locking it for further
// writes. The dataset is guaranteed to remain uncommitted on failure.
func (h *DatasetHandle) Commit(ctx context.Context) error {
	path := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	body := api.DatasetPatchSpec{Commit: true}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Delete a dataset. Note that this action is not reversible.
func (h *DatasetHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/datasets", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (c *Client) SearchDatasets(
	ctx context.Context,
	searchOptions api.DatasetSearchOptions,
	page int,
) ([]api.Dataset, error) {
	query := url.Values{"page": {strconv.Itoa(page)}}
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/datasets/search", query, searchOptions)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body []api.Dataset
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
}
