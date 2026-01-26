// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"context"
	"time"

	"github.com/matthewpi/nxhttp/httpheader"
	"github.com/matthewpi/nxretry"
)

// options for a [Client].
type options struct {
	//
	// Request defaults
	//

	// defaultHeaders for all http requests.
	defaultHeaders map[httpheader.Key]string

	//
	// nxretry
	//

	// backoff is used to control the delay between attempts.
	backoff nxretry.Backoff

	// maxAttempts is the maximum number of attempts that can occur. If set to 0
	// the maximum number of attempts will be unlimited.
	maxAttempts uint

	// minRetryAfter sets the minimum value of a "Retry-After" response header
	// that we will respect from a server.
	//
	// If the duration of the Retry-After is less than what is configured here,
	// the Retry-After is ignored and the next interval from the backoff is
	// used.
	minRetryAfter time.Duration

	// maxRetryAfter sets the maximum value of a "Retry-After" response header
	// that we will respect from a server.
	//
	// If the duration of the Retry-After exceeds this amount, it will be
	// truncated (`min(Retry-After, maxRetryAfter)`).
	maxRetryAfter time.Duration

	//
	// other options
	//

	// onError .
	// TODO: document
	onError ErrorFunc

	// onErrorResponse .
	// TODO: document
	onErrorResponse ErrorResponseFunc
}

// newOptions creates a new [options] instance with any defaults.
func newOptions() *options {
	return &options{
		maxAttempts:   3,
		minRetryAfter: 3 * time.Second,
		maxRetryAfter: 30 * time.Second,
	}
}

// setDefaults sets the defaults for [options]. This is expected to be used
// after applying any [Option] to ensure any required options are set.
func (o *options) setDefaults() {
	if o.backoff == nil {
		o.backoff = &nxretry.Exponential{
			Factor: 2,
			Min:    1 * time.Second,
			Max:    5 * time.Second,
		}
	}

	if o.onError == nil {
		o.onError = func(_ context.Context, err error) error { return err }
	}
}

// Option for an [Client].
type Option interface {
	// apply applies the [Option] to an [options] instance.
	apply(*options)
}

// OptionFunc type is an adapter to allow the use of ordinary functions as an
// [Option]. If f is a function with the appropriate signature, `OptionFunc(f)`
// is an [Option] that calls f.
type OptionFunc func(o *options)

// Ensure that [OptionFunc] implements the [Option] interface.
var _ Option = (*OptionFunc)(nil)

// apply applies the [Option] to an [options] instance.
func (f OptionFunc) apply(o *options) { f(o) }

// Ensure that [OptionFunc] implements the [ClientOption] interface.
var _ ClientOption = (*OptionFunc)(nil)

// applyClient makes the [OptionFunc] type satisfy the [ClientOption] interface,
// so [Option] can be used as [ClientOption] but not the other way around.
func (f OptionFunc) applyClient(o *clientOptions) {}

//
// Request default options
//

// WithDefaultHeader sets a default header for all HTTP requests.
func WithDefaultHeader(k httpheader.Key, v string) OptionFunc {
	return func(o *options) {
		if o.defaultHeaders == nil {
			o.defaultHeaders = make(map[httpheader.Key]string)
		}
		o.defaultHeaders[k] = v
	}
}

// WithDefaultHeaders sets the default headers for all HTTP requests.
func WithDefaultHeaders(h map[httpheader.Key]string) OptionFunc {
	return func(o *options) {
		// TODO: merge into existing map or just replace it like we currently do?
		o.defaultHeaders = h
	}
}

//
// nxretry options
//

// MaxAttempts sets the maximum number of attempts that can occur. If set to 0
// the maximum number of attempts will be unlimited.
func MaxAttempts(maxAttempts uint) OptionFunc {
	return func(o *options) { o.maxAttempts = maxAttempts }
}

// MinRetryAfter .
// TODO: document
func MinRetryAfter(d time.Duration) OptionFunc {
	return func(o *options) { o.minRetryAfter = d }
}

// MaxRetryAfter .
// TODO: document
func MaxRetryAfter(d time.Duration) OptionFunc {
	return func(o *options) { o.maxRetryAfter = d }
}

//
// other options
//

// ErrorFunc .
// TODO: document
type ErrorFunc func(context.Context, error) error

// OnError .
// TODO: document
func OnError(fn ErrorFunc) OptionFunc {
	return func(o *options) { o.onError = fn }
}

// ErrorResponseFunc .
// TODO: document
type ErrorResponseFunc func(context.Context, *Response) error

// OnErrorResponse .
// TODO: document
//
// NOTE: consuming the response body in this function prevents it from being
// used by the main response handler. If you must consume the response body
// here, consider buffering and making it available by resetting `Body` on
// the response.
func OnErrorResponse(fn ErrorResponseFunc) OptionFunc {
	return func(o *options) { o.onErrorResponse = fn }
}
