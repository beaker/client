package client

import (
	"context"
	"net/http"
	"path"

	"github.com/pkg/errors"

	"github.com/beaker/client/api"
)

// MetricAlerts returns all enabled metric alerts. Intended for use by the alerts service only.
func (c *Client) MetricAlerts(ctx context.Context) (*api.MetricAlerts, error) {
	path := "/api/v3/alerts/"
	resp, err := c.sendRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var alerts api.MetricAlerts
	if err := parseResponse(resp, &alerts); err != nil {
		return nil, err
	}

	return &alerts, nil
}

// AlertHandle provides operations on an alert.
type AlertHandle struct {
	client *Client
	id     string
}

// MetricAlert gets a handle for a metric alert by ID.
func (c *Client) MetricAlert(ctx context.Context, id string) (*AlertHandle, error) {
	id, err := c.resolveRef(ctx, "/api/v3/alerts", id)
	if err != nil {
		return nil, errors.WithMessage(err, "could not resolve alert id "+id)
	}

	return &AlertHandle{client: c, id: id}, nil
}

// Trigger sets the alert's Triggered status
func (h *AlertHandle) Trigger(ctx context.Context, trigger bool) error {
	path := path.Join("/api/v3/alerts/", h.id)
	body := api.MetricAlertPatchSpec{
		Triggered: &trigger,
	}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// Enable sets the alert's Enabled status
func (h *AlertHandle) Enable(ctx context.Context, enable bool) error {
	path := path.Join("/api/v3/alerts/", h.id)
	body := api.MetricAlertPatchSpec{
		Enabled: &enable,
	}
	resp, err := h.client.sendRequest(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
