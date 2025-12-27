package apptest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

// Client is a simple HTTP client for testing purposes.
type Client struct {
	httpCli *http.Client
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpCli: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) do(t *testing.T, method, url, contentType string, data []byte) (string, int) {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("could not create a HTTP request: %v", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := c.httpCli.Do(req)

	if err != nil {
		t.Fatalf("could not send HTTP request: %v", err)
	}
	body := readAllAndClose(t, res.Body)

	return body, res.StatusCode
}

func (c *Client) CloseConnections() {
	c.httpCli.CloseIdleConnections()
}

func (c *Client) Get(t *testing.T, url string) (string, int) {
	t.Helper()
	return c.do(t, http.MethodGet, url, "", nil)
}

func (c *Client) Post(t *testing.T, url string, data []byte) (string, int) {
	t.Helper()
	return c.do(t, http.MethodPost, url, "application/json", data)
}

func (c *Client) Delete(t *testing.T, url string) (string, int) {
	t.Helper()
	return c.do(t, http.MethodDelete, url, "", nil)
}

func readAllAndClose(t *testing.T, rc io.ReadCloser) string {
	t.Helper()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}
	if cerr := rc.Close(); cerr != nil {
		t.Fatalf("could not close response body: %v", cerr)
	}
	return string(b)
}
