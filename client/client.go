package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/goware/urlx"
	retryable "github.com/hashicorp/go-retryablehttp"

	"github.com/beaker/client/api"
)

var (
	// If version is empty the client will not send a version header.
	// If it is not empty, the client will send a version header and the service
	// will respond with an error if the client version is out of date.
	// The CLI sets version at link time.
	version = ""
)

// Client is a Beaker HTTP client.
type Client struct {
	baseURL   url.URL
	userToken string

	userAgent string
	// If set, then HTTPResponseHook will be invoked after every HTTP response
	// arrives. This can be used by users of this client to implement diagnostics,
	// such as request logging.
	HTTPResponseHook HTTPResponseHook
}

// Optional additions while constructing the client.
type Option interface {
	apply(*Client) error
}

type optionFunc func(*Client) error

func (o optionFunc) apply(c *Client) error {
	return o(c)
}

// WithUserAgent sets the useragent header for all requests issued by the client.
// If not specified, the executable's name is used as the default user-agent.
func WithUserAgent(userAgent string) Option {
	return optionFunc(func(c *Client) error {
		c.userAgent = userAgent
		return nil
	})
}

// HTTPResponseHook will be given an HTTP response, and the duration that the request took.
// When inspecting the response, don't read or close the response body, as that will affect
// the client behavior.
type HTTPResponseHook func(resp *http.Response, duration time.Duration)

// NewClient creates a new Beaker client bound to a single user.
func NewClient(address string, userToken string, opts ...Option) (*Client, error) {
	u, err := urlx.ParseWithDefaultScheme(address, "https")
	if err != nil {
		return nil, err
	}

	if u.Path != "" || u.Opaque != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return nil, errors.New("address must be base server address in the form [scheme://]host[:port]")
	}

	client := &Client{
		baseURL:   *u,
		userToken: userToken,
	}

	for _, opt := range opts {
		if err := opt.apply(client); err != nil {
			return nil, err
		}
	}

	if client.userAgent == "" {
		exec, err := os.Executable()
		if err != nil {
			exec = os.Args[0]
		}
		client.userAgent = path.Base(exec)
	}

	return client, nil
}

func newRetryableClient(httpClient *http.Client, httpResponseHook HTTPResponseHook) *retryable.Client {
	rc := &retryable.Client{
		HTTPClient:   httpClient,
		Logger:       &errorLogger{Logger: log.New(os.Stderr, "", log.LstdFlags)},
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 30 * time.Second,
		RetryMax:     9,
		CheckRetry:   retryable.DefaultRetryPolicy,
		Backoff:      exponentialJitterBackoff,
		ErrorHandler: retryable.PassthroughErrorHandler,
	}

	if httpResponseHook != nil {
		th := &timingHook{responseHook: httpResponseHook}
		rc.RequestLogHook = th.RequestLogHook
		rc.ResponseLogHook = th.ResponseLogHook
	}

	return rc
}

type errorLogger struct {
	Logger *log.Logger
}

func (l *errorLogger) Printf(template string, args ...interface{}) {
	if strings.HasPrefix(template, "[ERR]") {
		l.Logger.Printf(template, args...)
	}
}

// Address returns a client's host and port.
func (c *Client) Address() string {
	return c.baseURL.String()
}

func (c *Client) sendRetryableRequest(
	ctx context.Context,
	method string,
	path string,
	query url.Values,
	body interface{},
) (*http.Response, error) {
	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := c.newRetryableRequest(method, path, query, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return newRetryableClient(&http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: copyRedirectHeader,
	}, c.HTTPResponseHook).Do(req.WithContext(ctx))
}

func (c *Client) newRequest(
	method string,
	path string,
	query url.Values,
	body io.Reader,
) (*http.Request, error) {
	u := c.baseURL.ResolveReference(&url.URL{Path: path, RawQuery: query.Encode()})
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if version != "" {
		req.Header.Set(api.HeaderVersion, version)
	}
	if len(c.userToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.userToken)
	}

	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

func (c *Client) newRetryableRequest(
	method string,
	path string,
	query url.Values,
	body io.Reader,
) (*retryable.Request, error) {
	u := c.baseURL.ResolveReference(&url.URL{Path: path, RawQuery: query.Encode()})
	req, err := retryable.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	if version != "" {
		req.Header.Set(api.HeaderVersion, version)
	}
	if len(c.userToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.userToken)
	}

	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

func copyRedirectHeader(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}
	for key, val := range via[0].Header {
		req.Header[key] = val
	}
	return nil
}

// errorFromResponse creates an error from an HTTP response, or nil on success.
func errorFromResponse(resp *http.Response) error {
	// Anything less than 400 isn't an error, so don't produce one.
	if resp.StatusCode < 400 {
		return nil
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var apiErr api.Error
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return fmt.Errorf("failed to parse response: %s", string(bytes))
	}

	return apiErr
}

// responseValue parses the response body and stores the result in the given value.
// The value parameter should be a pointer to the desired structure.
func parseResponse(resp *http.Response, value interface{}) error {
	if err := errorFromResponse(resp); err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, value)
}

// safeClose closes an object while safely handling nils.
func safeClose(closer io.Closer) {
	if closer == nil {
		return
	}
	_ = closer.Close()
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// exponentialJitterBackoff implements exponential backoff with full jitter as described here:
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
func exponentialJitterBackoff(
	minDuration, maxDuration time.Duration,
	attempt int,
	resp *http.Response,
) time.Duration {
	min := float64(minDuration)
	max := float64(maxDuration)

	backoff := min + math.Min(max-min, min*math.Exp2(float64(attempt)))*random.Float64()
	return time.Duration(backoff)
}

type timingHook struct {
	start        time.Time
	responseHook HTTPResponseHook
}

func (th *timingHook) RequestLogHook(logger retryable.Logger, req *http.Request, attemptNum int) {
	th.start = time.Now()
}

func (th *timingHook) ResponseLogHook(logger retryable.Logger, resp *http.Response) {
	duration := time.Since(th.start)
	th.responseHook(resp, duration)
}
