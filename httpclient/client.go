package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client with retry and timeout capabilities
type Client struct {
	httpClient *http.Client
	config     Config
}

// Config holds HTTP client configuration
type Config struct {
	Timeout        time.Duration     `json:"timeout" yaml:"timeout"`
	Retries        int               `json:"retries" yaml:"retries"`
	RetryDelay     time.Duration     `json:"retry_delay" yaml:"retry_delay"`
	UserAgent      string            `json:"user_agent" yaml:"user_agent"`
	DefaultHeaders map[string]string `json:"default_headers" yaml:"default_headers"`
}

// Response represents an HTTP response
type Response struct {
	*http.Response
	Body []byte
}

// RequestOption allows customizing individual requests
type RequestOption func(*http.Request)

// New creates a new HTTP client with default configuration
func New() *Client {
	return NewWithConfig(Config{
		Timeout:    30 * time.Second,
		Retries:    3,
		RetryDelay: 1 * time.Second,
		UserAgent:  "gopackkit-httpclient/1.0",
	})
}

// NewWithConfig creates a new HTTP client with custom configuration
func NewWithConfig(config Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// Get performs a GET request
func (c *Client) Get(url string, opts ...RequestOption) (*Response, error) {
	return c.Do(http.MethodGet, url, nil, opts...)
}

// Post performs a POST request
func (c *Client) Post(url string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.Do(http.MethodPost, url, body, opts...)
}

// Put performs a PUT request
func (c *Client) Put(url string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.Do(http.MethodPut, url, body, opts...)
}

// Delete performs a DELETE request
func (c *Client) Delete(url string, opts ...RequestOption) (*Response, error) {
	return c.Do(http.MethodDelete, url, nil, opts...)
}

// Do performs an HTTP request with retry logic
func (c *Client) Do(method, url string, body interface{}, opts ...RequestOption) (*Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		resp, err := c.doRequest(method, url, body, opts...)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on the last attempt
		if attempt < c.config.Retries {
			time.Sleep(c.config.RetryDelay)
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.Retries+1, lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(method, url string, body interface{}, opts ...RequestOption) (*Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for key, value := range c.config.DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Set User-Agent
	if c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}

	// Set Content-Type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply request options
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error or handle as appropriate
			_ = closeErr
		}
	}()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		Response: resp,
		Body:     respBody,
	}, nil
}

// WithHeader adds a header to the request
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// WithHeaders adds multiple headers to the request
func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
}

// WithContext adds a context to the request
func WithContext(ctx context.Context) RequestOption {
	return func(req *http.Request) {
		*req = *req.WithContext(ctx)
	}
}

// WithAuth adds basic authentication to the request
func WithAuth(username, password string) RequestOption {
	return func(req *http.Request) {
		req.SetBasicAuth(username, password)
	}
}

// WithBearerToken adds a bearer token to the request
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// JSON unmarshals the response body into the provided interface
func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// String returns the response body as a string
func (r *Response) String() string {
	return string(r.Body)
}
