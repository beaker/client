package client

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/beaker/client/api"
)

// CreateImage creates a new image with an optional name.
func (c *Client) CreateImage(
	ctx context.Context,
	spec api.ImageSpec,
	name string,
) (*ImageHandle, error) {
	query := url.Values{}
	if name != "" {
		query.Set("name", name)
	}

	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/images", query, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.Image
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}

	return &ImageHandle{client: c, ref: body.ID}, nil
}

// Image gets a handle for an image by name or ID. The reference is not resolved.
func (c *Client) Image(reference string) *ImageHandle {
	return &ImageHandle{client: c, ref: reference}
}

// ImageHandle provides operations on an image.
type ImageHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *ImageHandle) Ref() string {
	return h.ref
}

// Get retrieves an image's details.
func (h *ImageHandle) Get(ctx context.Context) (*api.Image, error) {
	uri := path.Join("/api/v3/images", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, uri, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.Image
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// Repository returns information required to access an image through Docker.
func (h *ImageHandle) Repository(
	ctx context.Context,
	upload bool,
) (*api.ImageRepository, error) {
	path := path.Join("/api/v3/images", url.PathEscape(h.ref), "repository")
	query := url.Values{"upload": {strconv.FormatBool(upload)}}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body api.ImageRepository
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

// SetName sets an image's name.
func (h *ImageHandle) SetName(ctx context.Context, name string) error {
	path := path.Join("/api/v3/images", url.PathEscape(h.ref))
	body := api.ImagePatchSpec{Name: &name}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// SetDescription sets an image's description.
func (h *ImageHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/images", url.PathEscape(h.ref))
	body := api.ImagePatchSpec{Description: &description}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Commit finalizes an image, unblocking usage and locking it for further
// writes. The image is guaranteed to remain uncommitted on failure.
func (h *ImageHandle) Commit(ctx context.Context) error {
	path := path.Join("/api/v3/images", url.PathEscape(h.ref))
	body := api.ImagePatchSpec{Commit: true}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Delete an image. Note that this action is not reversible.
func (h *ImageHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/images", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (c *Client) SearchImages(
	ctx context.Context,
	searchOptions api.ImageSearchOptions,
	page int,
) ([]api.Image, error) {
	query := url.Values{"page": {strconv.Itoa(page)}}
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/images/search", query, searchOptions)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body []api.Image
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
}
