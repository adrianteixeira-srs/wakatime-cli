package api

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/wakatime/wakatime-cli/pkg/log"
)

const (
	// BaseURL is the base url of the wakatime api.
	BaseURL = "https://services-dev.roofhub.pro/copilot-metrics"
	// BaseIPAddrv4 is the base ip address v4 of the wakatime api.
	BaseIPAddrv4 = "143.244.210.202"
	// BaseIPAddrv6 is the base ip address v6 of the wakatime api.
	BaseIPAddrv6 = "2604:a880:4:1d0::2a7:b000"
	// DefaultTimeoutSecs is the default timeout used for requests to the wakatime api.
	DefaultTimeoutSecs = 120
	// MaxRetries is the maximum number of retries for requests to the wakatime api.
)

// Client communicates with the wakatime api.
type Client struct {
	baseURL string
	client  *http.Client
	// doFunc allows api client options to manipulate request/response handling.
	// default function will be set in constructor.
	//
	// wrapping by api options should be performed as follows:
	//
	//	next := c.doFunc
	//	c.doFunc = func(c *Client, req *http.Request) (*http.Response, error) {
	//		// do something
	//		resp, err := next(c, req)
	//		// do more
	//		return resp, err
	//	}
	doFunc func(c *Client, req *http.Request) (*http.Response, error)
}

// NewClient creates a new Client. Any number of Options can be provided.
func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL: baseURL,
		client: &http.Client{
			Transport: NewTransport(),
		},
		doFunc: func(c *Client, req *http.Request) (*http.Response, error) {
			req.Header.Set("Accept", "application/json")
			return c.client.Do(req)
		},
	}

	for _, option := range opts {
		option(c)
	}

	return c
}

// Do executes c.doFunc(), logs request/response details, and returns the response or error.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	logger := log.Extract(ctx)
	// Monta comando curl equivalente para debug/reprodução
	curlCmd := []string{"curl", "-X", req.Method, "'" + req.URL.String() + "'"}
	for k, v := range req.Header {
		for _, vv := range v {
			curlCmd = append(curlCmd, "-H", "'"+k+": "+vv+"'")
		}
	}
	if req.Body != nil && req.GetBody != nil {
		rc, _ := req.GetBody()
		if rc != nil {
			bodyCopy, _ := io.ReadAll(rc)
			rc.Close()
			if len(bodyCopy) > 0 {
				curlCmd = append(curlCmd, "--data-raw", "'"+strings.ReplaceAll(string(bodyCopy), "'", "\\'")+"'")
			}
		}
	}
	logger.Infof("[API] cURL: %s", strings.Join(curlCmd, " "))

	// NÃO sobrescrever Authorization aqui! O valor correto deve ser setado antes da chamada a Do().

	// Log request details
	logger.Infof("[API] Request: %s %s", req.Method, req.URL.String())
	logger.Debugf("[API] Request Headers: %v", req.Header)
	// Log Authorization header explicitly for debugging
	if auth := req.Header.Get("Authorization"); auth != "" {
		logger.Infof("[API] Request Authorization: %s", auth)
	} else {
		logger.Infof("[API] Request Authorization: <none>")
	}
	if req.Body != nil {
		// Try to read and log the body (if possible)
		var bodyCopy []byte
		if req.GetBody != nil {
			rc, _ := req.GetBody()
			if rc != nil {
				bodyCopy, _ = io.ReadAll(rc)
				rc.Close()
			}
		}
		if len(bodyCopy) > 0 {
			logger.Debugf("[API] Request Body: %s", string(bodyCopy))
		}
	}

	resp, err := c.doFunc(c, req)
	if err != nil {
		logger.Errorf("[API] Request error: %v", err)
		return nil, err
	}

	// Log response details
	logger.Infof("[API] Response: %s %s -> %d", req.Method, req.URL.String(), resp.StatusCode)
	logger.Debugf("[API] Response Headers: %v", resp.Header)
	if resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Debugf("[API] Response Body: %s", string(bodyBytes))
		// Restore body for further use
		resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	return resp, nil
}
