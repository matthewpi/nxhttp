// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/matthewpi/nxhttp/httpheader"
)

// isTimeout checks if err was the result of a timeout.
func isTimeout(err error) bool {
	// Check if the error is a context timeout, many stdlib errors are comparable
	// to [context.DeadlineExceeded].
	//
	// Examples:
	// - https://github.com/golang/go/blob/go1.24.5/src/net/net.go#L627
	// - https://github.com/golang/go/blob/go1.24.5/src/net/http/transport.go#L2713
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Just to ensure we catch all the edge-cases, also check the `Timeout()`
	// method if available.
	var tErr interface{ Timeout() bool }
	return errors.As(err, &tErr) && tErr.Timeout()
}

// ContentError indicates an HTTP response contained an unexpected `Content-*`
// header.
type ContentError struct {
	// Header that triggered the error.
	Header httpheader.Key
	// Value of the header in the response.
	Value string
	// Allowed or expected values for the header.
	Allowed []string
}

var (
	_ error          = ContentError{}
	_ slog.LogValuer = ContentError{}
)

// NewContentError returns a new [ContentError].
func NewContentError(key httpheader.Key, value string, allowed ...string) ContentError {
	if len(allowed) == 0 {
		allowed = []string{""}
	}
	return ContentError{
		Header:  key,
		Value:   value,
		Allowed: allowed,
	}
}

// Error returns an error message and satisfies the [error] interface.
func (e ContentError) Error() string {
	// If only one value was allowed, print it in a nicer format instead of
	// as a slice with a single item.
	if len(e.Allowed) == 1 {
		return fmt.Sprintf("nxhttp: expected '%s' header to match '%s', but got '%s' instead", e.Header, e.Allowed[0], e.Value)
	}
	return fmt.Sprintf("nxhttp: expected '%s' header to match one of %v, but got '%s' instead", e.Header, e.Allowed, e.Value)
}

// LogValue returns an [slog.Value] and satisfies the [slog.LogValuer] interface.
func (e ContentError) LogValue() slog.Value {
	// TODO: better log value.
	return slog.GroupValue(
		slog.String("message", e.Error()),
	)
}

// RequestError is returned if the request fails to be done, i.e. the server is
// never reached.
type RequestError struct {
	// err from the Request.
	err error
}

var (
	_ error          = RequestError{}
	_ slog.LogValuer = RequestError{}
)

// Error returns an error message and satisfies the [error] interface.
func (e RequestError) Error() string {
	return "nxhttp request failed: " + e.err.Error()
}

// LogValue returns an [slog.Value] and satisfies the [slog.LogValuer] interface.
func (e RequestError) LogValue() slog.Value {
	return slog.GroupValue(slog.Any("err", e.err))
}

// Unwrap returns the underlying [error] that caused the [RequestError].
func (e RequestError) Unwrap() error {
	return e.err
}

// StatusError indicates an HTTP request failure with a status code from a
// remote HTTP server.
type StatusError struct {
	// Data from the HTTP response.
	Data []byte

	// StatusCode is the status code from the [Response] that caused this
	// error.
	StatusCode int

	// Expected is the status code we expected to see in the [Response] but
	// we got [StatusCode] instead.
	Expected int
}

var (
	_ error          = StatusError{}
	_ slog.LogValuer = StatusError{}
)

// NewStatusError returns a new [StatusError].
func NewStatusError(res *Response, expected int) StatusError {
	e := StatusError{Expected: expected}
	if res == nil {
		return e
	}
	e.StatusCode = res.StatusCode
	if res.Body != nil {
		if b, err := io.ReadAll(http.MaxBytesReader(nil, res.Body, 4*1024)); err == nil {
			e.Data = bytes.TrimSpace(b)
		}
		_ = res.Body.Close()
	}
	return e
}

// AsStatusError attempts to recursively map err to a [StatusError], even if
// it has been wrapped.
//
// Internally, this uses [errors.As] and is provided as a convenience.
func AsStatusError(err error) (StatusError, bool) {
	var sErr StatusError
	if err == nil || !errors.As(err, &sErr) {
		return sErr, false
	}
	return sErr, true
}

// Error returns an error message and satisfies the [error] interface.
func (e StatusError) Error() string {
	return fmt.Sprintf("nxhttp: expected %d status code, but got %d (%q)", e.Expected, e.StatusCode, e.Data)
}

// LogValue returns an [slog.Value] and satisfies the [slog.LogValuer] interface.
func (e StatusError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", e.Error()),
		slog.GroupAttrs("status_code",
			slog.Int("got", e.StatusCode),
			slog.Int("expected", e.Expected),
		),
	)
}
