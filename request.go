// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"context"
	"io"
	"net/http"

	"github.com/matthewpi/nxhttp/httpheader"
)

// Request .
// TODO: document
type Request struct {
	*http.Request

	// body for the request.
	body BodyFunc
}

var _ io.WriterTo = (*Request)(nil)

// NewRequest returns a new [Request] given a method, URL, and optional body.
func NewRequest(ctx context.Context, method, url string, body any) (*Request, error) {
	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req := &Request{Request: httpReq}
	if err := req.SetBody(body); err != nil {
		return nil, err
	}
	return req, nil
}

// SetBody sets the body on the [Request].
func (r *Request) SetBody(v any) (err error) {
	r.body, r.ContentLength, err = GetBody(v)
	if r.body == nil {
		r.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	} else {
		r.GetBody = r.body
	}
	return err
}

// WithContext returns wrapped Request with a shallow copy of r with its context
// changed to ctx. The provided ctx must be non-nil.
func (r *Request) WithContext(ctx context.Context) *Request {
	return &Request{
		Request: r.Request.WithContext(ctx),
		body:    r.body,
	}
}

// AddHeader is like [http.Header.Add], but the key must already be in
// [httpheader.Canonicalize] form.
func (r *Request) AddHeader(key httpheader.Key, value string) {
	httpheader.Add(r.Header, key, value)
}

// SetHeader is like [http.Header.Set], but the key must already be in
// [httpheader.Canonicalize] form.
func (r *Request) SetHeader(key httpheader.Key, value string) {
	httpheader.Set(r.Header, key, value)
}

// DelHeader is like [http.Header.Del], but the key must already be in
// [httpheader.Canonicalize] form.
func (r *Request) DelHeader(key httpheader.Key) {
	httpheader.Del(r.Header, key)
}

// WriteTo implements the [io.WriterTo] interface.
func (r *Request) WriteTo(w io.Writer) (int64, error) {
	body, err := r.body()
	if err != nil {
		return 0, err
	}
	defer body.Close()
	return io.Copy(w, body)
}
