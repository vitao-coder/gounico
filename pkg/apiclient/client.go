package apiclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	cl *http.Client
}

func (c *Client) Get(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return do(ctx, client, req)
}

func (c *Client) Post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return do(ctx, client, req)
}

func (c *Client) PostForm(ctx context.Context, client *http.Client, url string, data url.Values) (*http.Response, error) {
	return c.Post(ctx, client, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}
	return resp, err
}
