// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import "net/http"

// requestOptions represent the options for a [Request].
type requestOptions struct {
	transport    func(t *http.Transport)
	roundTripper func(http.RoundTripper) http.RoundTripper
}

// RequestOption for an [Request].
type RequestOption interface {
	// apply applies the [RequestOption] to a [requestOptions] instance.
	apply(*requestOptions)
}

// RequestOptionFunc type is an adapter to allow the use of ordinary functions as
// a [RequestOption]. If f is a function with the appropriate signature,
// `RequestOptionFunc(f)` is a [RequestOption] that calls f.
type RequestOptionFunc func(o *requestOptions)

// Ensure that [RequestOptionFunc] implements the [RequestOption] interface.
var _ RequestOption = (*RequestOptionFunc)(nil)

// apply applies the [RequestOption] to a [requestOptions] instance.
func (f RequestOptionFunc) apply(o *requestOptions) { f(o) }

// WithRequestTransport sets the underlying [*http.Transport] that will be used
// for an individual request.
func WithRequestTransport(fn func(t *http.Transport)) RequestOptionFunc {
	return func(o *requestOptions) { o.transport = fn }
}

// WithRequestRoundTripper sets the underlying [http.RoundTripper] that will be
// used for an individual request.
func WithRequestRoundTripper(fn func(http.RoundTripper) http.RoundTripper) RequestOptionFunc {
	return func(o *requestOptions) { o.roundTripper = fn }
}
