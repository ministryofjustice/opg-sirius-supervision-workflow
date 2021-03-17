package sirius

import (
	"context"
	"io"
	"net/http"
)

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}, nil
}

type Client struct {
	http    *http.Client
	baseURL string
}

func (c *Client) newRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")
	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)

	return req, err
}

// func (c *Client) Authenticate(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, c.url("/auth"), http.StatusFound)
// }

// // func (c *Client) url(path string) string {
// // 	partial, _ := url.Parse(path)

// // 	return c.baseURL.ResolveReference(partial).String()
// // }

// func (e ClientError) Error() string {
// 	return string(e)
// }
