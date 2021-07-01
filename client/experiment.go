package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/beaker/client/api"
)

// Experiment gets a handle for an experiment by name or ID. The reference is not resolved.
func (c *Client) Experiment(reference string) *ExperimentHandle {
	return &ExperimentHandle{client: c, ref: reference}
}

// ExperimentHandle provides operations on an experiment.
type ExperimentHandle struct {
	client *Client
	ref    string
}

// Ref returns the name or ID with which a handle was created.
func (h *ExperimentHandle) Ref() string {
	return h.ref
}

// Get retrieves an experiment's details, including a summary of contained tasks.
func (h *ExperimentHandle) Get(ctx context.Context) (*api.Experiment, error) {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var experiment api.Experiment
	if err := parseResponse(resp, &experiment); err != nil {
		return nil, err
	}

	return &experiment, nil
}

// Groups gets the ID of each group that the experiment belongs to.
func (h *ExperimentHandle) Groups(ctx context.Context) ([]string, error) {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref), "groups")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var groups []string
	if err := parseResponse(resp, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// SetName sets an experiment's name.
func (h *ExperimentHandle) SetName(ctx context.Context, name string) error {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref))
	body := api.ExperimentPatchSpec{Name: &name}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// SetDescription sets an experiment's description
func (h *ExperimentHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref))
	body := api.ExperimentPatchSpec{Description: &description}
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Spec gets the experiment specification.
// Default format is YAML. JSON is available by setting json=true.
func (h *ExperimentHandle) Spec(ctx context.Context, version string, json bool) (io.ReadCloser, error) {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref), "spec")
	req, err := h.client.newRetryableRequest("GET", path, url.Values{
		"version": []string{version},
	}, nil)
	if err != nil {
		return nil, err
	}
	if json {
		req.Header.Set("Accept", "application/json")
	}

	resp, err := newRetryableClient(&http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}, h.client.HTTPResponseHook).Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err := errorFromResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

// Resume retries failed or stopped tasks within a previously run experiment.
func (h *ExperimentHandle) Resume(ctx context.Context) error {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref), "/resume")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Stop cancels all uncompleted tasks for an experiment. If the experiment has
// already completed, this succeeds without effect.
func (h *ExperimentHandle) Stop(ctx context.Context) error {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref), "stop")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPut, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Delete an experiment. This action is not reversible.
func (h *ExperimentHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Tasks of the experiment
func (h *ExperimentHandle) Tasks(ctx context.Context) ([]api.Task, error) {
	path := path.Join("/api/v3/experiments", url.PathEscape(h.ref), "tasks")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)
	var tasks []api.Task
	if err := parseResponse(resp, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Client) SearchExperiments(
	ctx context.Context,
	searchOptions api.ExperimentSearchOptions,
	page int,
) ([]api.Experiment, error) {
	query := url.Values{"page": {strconv.Itoa(page)}}
	resp, err := c.sendRetryableRequest(ctx, http.MethodPost, "/api/v3/experiments/search", query, searchOptions)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var body []api.Experiment
	if err := parseResponse(resp, &body); err != nil {
		return nil, err
	}

	return body, nil
}
