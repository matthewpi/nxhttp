// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"io"
	"net/http"
	"time"

	"github.com/matthewpi/nxhttp/httpheader"
)

// Response represents the response from an HTTP request.
//
// The [Client] return Responses from servers once the response headers have
// been received. The response body is streamed on demand as the Body field
// is read.
type Response struct {
	*http.Response
}

var _ io.Closer = (*Response)(nil)

// Closes the body of r.
func (r *Response) Close() error {
	// If the body is somehow nil, do nothing.
	if r.Body == nil {
		return nil
	}

	// Close the response body, this will also automatically discard any
	// unread contents up to a limit due to `r.Body` being wrapped with
	// [discardReadCloser].
	return r.Body.Close()
}

// GetHeader is like [http.Header.Get], but the key must already be in
// [httpheader.Key] form.
func (r *Response) GetHeader(key httpheader.Key) string {
	return httpheader.Get(r.Header, key)
}

// retryAfter parses the "Retry-After" header and returns it as a
// [time.Duration].
//
// The value of the header could be either an HTTP Date (see [http.ParseTime])
// or a number of seconds.
func (r *Response) retryAfter() (time.Duration, error) {
	return httpheader.ParseRetryAfter(httpheader.Get(r.Header, httpheader.RetryAfter))
}

// discardReadCloser wraps an [io.ReadCloser], overriding it's Close method
// with one that discards any remaining content from the reader.
type discardReadCloser struct {
	io.ReadCloser
	// eof indicates whether the [io.ReadCloser] we are wrapping has ever
	// returned an [io.EOF] error.
	eof bool
}

// Read satisfies [io.Reader].
func (r *discardReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err == io.EOF {
		r.eof = true
	}
	return n, err
}

// Close satisfies [io.Closer].
func (r *discardReadCloser) Close() error {
	if !r.eof {
		// Discard the rest of the response body.
		discard(r.ReadCloser)
	}
	return r.ReadCloser.Close()
}

// discard copies a limited amount of data from an [io.Reader] to [io.Discard].
func discard(r io.Reader) {
	// We use an [io.LimitReader] here to protect against misbehaving (or even
	// malicious) servers from causing us to read a large amount of data.
	//
	// If the response body doesn't get entirely drained, it is unlikely
	// the connection will be reused, however that is still acceptable.
	_, _ = io.Copy(io.Discard, io.LimitReader(r, 16*1024))
}
