// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"net"
	"net/http"
	"time"
)

// clientOptions for constructing an [*http.Client].
type clientOptions struct {
	transport     *http.Transport
	roundTripper  func(http.RoundTripper) http.RoundTripper
	checkRedirect CheckRedirectFunc
	timeout       time.Duration
	cookieJar     http.CookieJar
}

// Client returns a newly constructed [*http.Client] using the options.
func (o *clientOptions) Client() *http.Client {
	if o.transport == nil {
		o.transport = defaultTransport()
	}

	// If the user configured a RoundTripper (not just a transport), use it
	// to wrap the [*http.Transport].
	//
	// This is commonly used for integrating OpenTelemetry (otelhttp) as an
	// example.
	var rt http.RoundTripper
	if o.roundTripper != nil {
		rt = o.roundTripper(o.transport)
	} else {
		rt = o.transport
	}

	return &http.Client{
		Transport:     rt,
		CheckRedirect: o.checkRedirect,
		Jar:           o.cookieJar,
		Timeout:       o.timeout,
	}
}

// ClientOption for an [*http.Client].
type ClientOption interface {
	// applyClient applies the [ClientOption] to a [clientOptions] instance.
	applyClient(*clientOptions)
}

// ClientOptionFunc type is an adapter to allow the use of ordinary functions as
// a [ClientOption]. If f is a function with the appropriate signature,
// `ClientOptionFunc(f)` is a [ClientOption] that calls f.
type ClientOptionFunc func(o *clientOptions)

// Ensure that [ClientOptionFunc] implements the [ClientOption] interface.
var _ ClientOption = (*ClientOptionFunc)(nil)

// applyClient applies the [ClientOption] to a [clientOptions] instance.
func (f ClientOptionFunc) applyClient(o *clientOptions) { f(o) }

// SetTransport sets the underlying [*http.Transport] that will be used.
//
// If you don't already have an [*http.Transport] available, consider using
// [WithTransport] instead.
//
// NOTE: SetTransport can be used alongside [WithTransport] as long as
// [WithTransport] is called after SetTransport. Calls to SetTransport
// completely replace the [*http.Transport] used by the client.
func SetTransport(t *http.Transport) ClientOptionFunc {
	return func(o *clientOptions) {
		o.transport = t
	}
}

// defaultTransport is like [http.DefaultTransport], except with better defaults.
func defaultTransport() *http.Transport {
	// Configure a [net.Dialer] with a lower default timeout.
	d := &net.Dialer{
		Timeout:   5 * time.Second, // [http.DefaultTransport] uses `30 * time.Second` by default.
		KeepAlive: 30 * time.Second,
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment, // Default
		DialContext:           d.DialContext,
		ForceAttemptHTTP2:     true,             // Default
		MaxIdleConns:          100,              // Default
		MaxIdleConnsPerHost:   4,                // [http.Transport] uses `2` by default if this field is set to `0`.
		IdleConnTimeout:       30 * time.Second, // [http.DefaultTransport] uses `90 * time.Second` by default
		TLSHandshakeTimeout:   5 * time.Second,  // [http.DefaultTransport] uses `10 * time.Second` by default
		ExpectContinueTimeout: 1 * time.Second,  // Default
		// Wait at most 10 seconds for the server to at least respond with
		// headers after we have fully written the request (including the body)
		// to the server.
		ResponseHeaderTimeout: 10 * time.Second,
	}
}

// WithTransport sets the underlying [*http.Transport] that will be used.
//
// If you already have or prefer to construct your own [*http.Transport], you
// can use [SetTransport] instead.
func WithTransport(fn func(t *http.Transport)) ClientOptionFunc {
	return func(o *clientOptions) {
		// If a transport doesn't exist, create a new one, using our defaults.
		if o.transport == nil {
			o.transport = defaultTransport()
		}

		// Call the user provided function so they can modify the transport.
		fn(o.transport)
	}
}

// WithRoundTripper sets the underlying [http.RoundTripper] that will be used
// by the [*http.Client].
//
// This option was designed to work with [WithTransport] to allow the easy
// configuration of an [*http.Transport] but also the ability to wrap it with
// a [http.RoundTripper]. Such as for use with OpenTelemetry via the [otelhttp]
// library as an example.
//
// [otelhttp]: https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
func WithRoundTripper(rtFunc func(http.RoundTripper) http.RoundTripper) ClientOptionFunc {
	return func(o *clientOptions) { o.roundTripper = rtFunc }
}

// CheckRedirectFunc is the type of the function for the `CheckRedirect` field
// on an [http.Client].
type CheckRedirectFunc func(req *http.Request, via []*http.Request) error

// CheckRedirect sets the policy for handling redirects.
func CheckRedirect(fn CheckRedirectFunc) ClientOptionFunc {
	return func(o *clientOptions) { o.checkRedirect = fn }
}

// UseLastResponse disables the following of redirects. This function is
// designed to be used with [CheckRedirect].
//
//	nxhttp.CheckRedirect(nxhttp.UseLastResponse)
func UseLastResponse(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

// WithCookieJar sets the [http.CookieJar] used by the [*http.Client].
func WithCookieJar(jar http.CookieJar) ClientOptionFunc {
	return func(o *clientOptions) { o.cookieJar = jar }
}

// WithTimeout sets the timeout used by the [*http.Client].
func WithTimeout(d time.Duration) ClientOptionFunc {
	return func(o *clientOptions) { o.timeout = d }
}
