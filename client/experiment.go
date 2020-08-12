package client

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/beaker/client/api"
)

// ExperimentHandle provides operations on an experiment.
type ExperimentHandle struct {
	client *Client
	id     string
}

// CreateExperiment creates a new experiment with an optional name.
func (c *Client) CreateExperiment(
	ctx context.Context,
	spec api.ExperimentSpec,
	name string,
	priority string,
) (*ExperimentHandle, error) {
	query := url.Values{"name": {name}, "priority": {priority}}
	resp, err := c.sendRequest(ctx, http.MethodPost, "/api/v3/experiments", query, spec)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var id string
	if err := parseResponse(resp, &id); err != nil {
		return nil, err
	}

	return &ExperimentHandle{client: c, id: id}, nil
}

// ResumeExperiment resumes a previously preempted experiment. The new experiment is
// created with an optional name if provided in experimentName. The experiment referenced
// in resumeFromReference should refer to the name or ID of the experiment to resume from.
// The experiment referenced must contain at least one 'preempted' task.
func (c *Client) ResumeExperiment(
	ctx context.Context,
	resumeFromReference string,
	experimentName string,
) (*ExperimentHandle, error) {
	id, err := c.resolveRef(ctx, "/api/v3/experiments", resumeFromReference)
	if err != nil {
		return nil, errors.WithMessage(err, "could not resolve experiment reference "+resumeFromReference)
	}
	query := url.Values{"name": {experimentName}}
	resp, err := c.sendRequest(ctx, http.MethodPost, path.Join("/api/v3/experiments", id, "/resume"), query, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var resumedExperimentID string
	if err := parseResponse(resp, &resumedExperimentID); err != nil {
		return nil, err
	}

	return &ExperimentHandle{client: c, id: resumedExperimentID}, nil
}

// Experiment gets a handle for an experiment by name or ID. The returned handle
// is guaranteed throughout its lifetime to refer to the same object, even if
// that object is later renamed.
func (c *Client) Experiment(ctx context.Context, reference string) (*ExperimentHandle, error) {
	id, err := c.resolveRef(ctx, "/api/v3/experiments", reference)
	if err != nil {
		return nil, errors.WithMessage(err, "could not resolve experiment reference "+reference)
	}

	return &ExperimentHandle{client: c, id: id}, nil
}

// ID returns an experiment's stable, unique ID.
func (h *ExperimentHandle) ID() string {
	return h.id
}

// Get retrieves an experiment's details, including a summary of contained tasks.
func (h *ExperimentHandle) Get(ctx context.Context) (*api.Experiment, error) {
	path := path.Join("/api/v3/experiments", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodGet, path, nil, nil)
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

// SetName sets an experiment's name.
func (h *ExperimentHandle) SetName(ctx context.Context, name string) error {
	path := path.Join("/api/v3/experiments", h.id)
	body := api.ExperimentPatchSpec{Name: &name}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// SetDescription sets an experiment's description
func (h *ExperimentHandle) SetDescription(ctx context.Context, description string) error {
	path := path.Join("/api/v3/experiments", h.id)
	body := api.ExperimentPatchSpec{Description: &description}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Stop cancels all uncompleted tasks for an experiment. If the experiment has
// already completed, this succeeds without effect.
func (h *ExperimentHandle) Stop(ctx context.Context) error {
	path := path.Join("/api/v3/experiments", h.id, "stop")
	resp, err := h.client.sendRequest(ctx, http.MethodPut, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Delete an experiment. This action is not reversible.
func (h *ExperimentHandle) Delete(ctx context.Context) error {
	path := path.Join("/api/v3/experiments", h.id)
	resp, err := h.client.sendRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

func (c *Client) SearchExperiments(
	ctx context.Context,
	searchOptions api.ExperimentSearchOptions,
	page int,
) ([]api.Experiment, error) {
	query := url.Values{"page": {strconv.Itoa(page)}}
	resp, err := c.sendRequest(ctx, http.MethodPost, "/api/v3/experiments/search", query, searchOptions)
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
