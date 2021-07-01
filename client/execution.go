package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/beaker/client/api"
)

// Execution gets a handle for an execution by ID. The id is not resolved.
func (c *Client) Execution(id string) *ExecutionHandle {
	return &ExecutionHandle{client: c, id: id}
}

// ExecutionHandle provides access to a single execution.
type ExecutionHandle struct {
	client *Client
	id     string
}

// Ref returns the name or ID with which a handle was created.
func (h *ExecutionHandle) Ref() string {
	return h.id
}

// Get retrieves an execution's details.
func (h *ExecutionHandle) Get(ctx context.Context) (*api.Execution, error) {
	path := path.Join("/api/v3/executions", url.PathEscape(h.id))
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var result api.Execution
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLogs gets all logs for a task. Logs are in the form:
// {RFC3339 nano timestamp} {message}\n
func (h *ExecutionHandle) GetLogs(ctx context.Context) (io.ReadCloser, error) {
	path := path.Join("/api/v3/executions", url.PathEscape(h.id), "logs")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := errorFromResponse(resp); err != nil {
		safeClose(resp.Body)
		return nil, err
	}
	return resp.Body, nil
}

// PutLogs uploads a log chunk. Since is the time of the first log message in the chunk.
func (h *ExecutionHandle) PutLogs(ctx context.Context, filename string, logs io.Reader) error {
	path := path.Join("/api/v3/executions", url.PathEscape(h.id), "logs", filename)
	req, err := h.client.newRequest(http.MethodPut, path, nil, logs)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}

// GetResults retrieves an execution's results.
func (h *ExecutionHandle) GetResults(ctx context.Context) (*api.ExecutionResults, error) {
	path := path.Join("/api/v3/executions", url.PathEscape(h.id), "results")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, err
	}
	defer safeClose(resp.Body)

	var results api.ExecutionResults
	if err := parseResponse(resp, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// PostStatus updates an execution's current status.
func (h *ExecutionHandle) PostStatus(ctx context.Context, status api.ExecStatusUpdate) error {
	path := path.Join("/api/v3/executions", url.PathEscape(h.id), "status")
	resp, err := h.client.sendRetryableRequest(ctx, http.MethodPost, path, nil, status)
	if err != nil {
		return err
	}
	defer safeClose(resp.Body)
	return errorFromResponse(resp)
}
