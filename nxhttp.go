// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

// Package nxhttp .
// TODO: document
package nxhttp

import (
	"fmt"
	"net/http"

	"github.com/matthewpi/nxretry"
)

// Client is an HTTP client.
type Client struct {
	*options

	client    *http.Client
	transport *http.Transport
}

// FromClient constructs a [Client] using an existing [*http.Client] and any
// provided [Option].
//
// If you are able, it is preferred for users to construct a [Client] using
// [NewClient]. You can use [ClientOption] to much more easily configure the
// underlying [*http.Client], especially if you need a custom [*http.Transport]
// and not just a custom [http.RoundTripper].
func FromClient(h *http.Client, opts ...Option) *Client {
	if h == nil {
		h = http.DefaultClient
	}
	o := newOptions()
	for _, opt := range opts {
		opt.apply(o)
	}
	o.setDefaults()
	c := &Client{
		options: o,
		client:  h,
	}
	if t, ok := h.Transport.(*http.Transport); ok {
		c.transport = t
	} else {
		c.transport = http.DefaultTransport.(*http.Transport)
	}
	return c
}

// NewClient constructs a [Client] using any provided [ClientOption].
func NewClient(opts ...ClientOption) *Client {
	co := &clientOptions{}
	o := newOptions()
	for _, clientOpt := range opts {
		if opt, ok := clientOpt.(Option); ok {
			opt.apply(o)
		} else {
			clientOpt.applyClient(co)
		}
	}
	o.setDefaults()
	return &Client{
		options:   o,
		client:    co.Client(),
		transport: co.transport,
	}
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(req *Request, opts ...RequestOption) (*Response, error) {
	ctx := req.Context()
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("nxhttp: context already has an error: %w", err)
	}

	// If configured, set any default headers on the request.
	if c.defaultHeaders != nil {
		for key := range c.defaultHeaders {
			if _, ok := req.Header[string(key)]; !ok {
				req.Header[string(key)] = []string{c.defaultHeaders[key]}
			}
		}
	}

	httpClient := c.client

	// Handle options for the request if present.
	if len(opts) > 0 {
		reqOpts := &requestOptions{}
		for _, opt := range opts {
			opt.apply(reqOpts)
		}

		var rt http.RoundTripper
		if reqOpts.transport != nil {
			t := c.transport.Clone()
			reqOpts.transport(t)
			rt = t
		}
		if reqOpts.roundTripper != nil {
			if rt == nil {
				rt = c.transport.Clone()
			}
			rt = reqOpts.roundTripper(rt)
		}

		// If the transport was overridden, create a new HTTP Client that
		// uses the transport.
		if rt != nil {
			httpClient = &http.Client{
				Transport:     rt,
				CheckRedirect: httpClient.CheckRedirect,
				Jar:           httpClient.Jar,
				Timeout:       httpClient.Timeout,
			}
		}
	}

	var (
		r     *Response
		doErr error
	)
	// Configure the retrier for the request.
	rty := nxretry.New(
		nxretry.MaxAttempts(c.maxAttempts),
		c.backoff,
	)
	for range rty.Next(ctx) {
		// Execute the request.
		r, doErr = doRequest(httpClient, req)
		if doErr != nil {
			// Allow the caller to process the error before we do.
			//
			// This can be used for logging or to transform errors such as
			// permanent to temporary or vice versa.
			doErr = c.onError(ctx, doErr)

			// Only retry here if the error is retryable. We don't want to keep
			// retrying a broken request such as one with a malformed URL, but
			// we do for a connection timeout (as an example).
			if doErr == nil || isTimeout(doErr) {
				continue
			}

			// Otherwise, treat the error as permanent and don't retry the request.
			break
		}

		// If we got a successful status code, return the response immediately
		// without any additional processing.
		if r.StatusCode >= http.StatusOK && r.StatusCode <= 299 {
			break
		}

		// If the caller provided a handler for an error response, call it
		// to determine if we should continue or not.
		if c.onErrorResponse != nil {
			doErr = c.onErrorResponse(ctx, r)
			if doErr != nil {
				break
			}
		}

		// Depending on the status code of the response, determine if the
		// request should be retried.
		//
		// TODO: add an option on the client to allow/deny additional codes.
		switch r.StatusCode {
		case http.StatusTooManyRequests:
		case http.StatusInternalServerError:
		case http.StatusBadGateway:
		case http.StatusServiceUnavailable:
		case http.StatusGatewayTimeout:
		default:
			// The request was either successful or we hit a fatal error, either way
			// we are done.
			break
		}

		// Get the duration we should wait from the Retry-After header.
		d, err := r.retryAfter()
		if err != nil {
			// TODO: do we want to do something about the Retry-After error?
			//
			// If there is an error it means a malformed Retry-After was sent
			// and we may want to let the user know that. That way they can
			// fix the server (if they control it) or inform the server operator
			// about the issue.
			continue
		}

		// Only override the retrier if the Retry-After was parsed and
		// is above our minimum, otherwise fallback to the standard
		// backoff.
		if d > c.minRetryAfter {
			// Ensure the duration does not exceed our configured maximum
			// if configured.
			if c.maxRetryAfter > 0 && d > c.maxRetryAfter {
				// Truncate the duration to our maximum value.
				d = c.maxRetryAfter
			}

			// We got a valid Retry-After from the server, use as next delay
			// for the retrier instead of whatever the normal backoff would
			// provide.
			rty.Override(d)
		}

		// Retry the request. If there was a Retry-After header in the
		// response, it will be respected. Otherwise, the [nxretry.Backoff]
		// that was configured will be used to determine the delay for the
		// next attempt.
		continue
	}

	// Return the response and error. It is very likely one of them is nil, but
	// that is to be expected.
	return r, doErr
}

// do wraps a [http.Client.Do] method.
func doRequest(c *http.Client, req *Request) (*Response, error) {
	// Check if our custom body type is set, while we end up using the
	// GetBody property anyways, the GetBody property will always be set,
	// but it will just return [http.NoBody] if no actual body is present.
	if req.body != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, RequestError{err: err}
		}
		req.Body = body
	}

	// Execute the request.
	res, err := c.Do(req.Request)
	if err != nil {
		return nil, RequestError{err: err}
	}

	// Avoid wrapping the response if it is somehow nil.
	if res == nil {
		// TODO: I don't ever believe that res should be nil, and to protect
		// downstream code from panicking do we want to error here instead of
		// passing through two nil return values?
		return nil, nil
	}

	// If we have a response with a body, wrap it with [discardReadCloser] so
	// when the body gets closed, we ensure its contents get read to completion
	// so the response can get reused for future requests.
	if res.Body != nil {
		res.Body = &discardReadCloser{ReadCloser: res.Body}
	}

	// Wrap the response.
	return &Response{Response: res}, nil
}
