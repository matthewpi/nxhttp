// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/matthewpi/nxhttp"
)

func Example() {
	_ = nxhttp.NewClient(
		// Client options
		nxhttp.WithTransport(func(t *http.Transport) {
			t.TLSClientConfig = &tls.Config{}
		}),
		nxhttp.WithTimeout(15*time.Second),
		// Regular options
		nxhttp.MaxAttempts(5),
		nxhttp.OnError(func(_ context.Context, err error) error { return err }),
	)
}

func ExampleSetTransport() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = &tls.Config{
		// Just as an example.
	}

	_ = nxhttp.NewClient(nxhttp.SetTransport(t))
}

func ExampleWithTransport() {
	_ = nxhttp.NewClient(
		nxhttp.WithTransport(func(t *http.Transport) {
			t.TLSClientConfig = &tls.Config{
				// Just as an example.
			}
		}),
	)
}

func ExampleWithRoundTripper() {
	_ = nxhttp.NewClient(
		nxhttp.WithRoundTripper(func(rt http.RoundTripper) http.RoundTripper {
			// NOTE: the following line is commented-out to avoid this library
			// depending on OpenTelemetry just for an example.

			// return otelhttp.NewTransport(rt)
			return rt
		}),
	)
}
