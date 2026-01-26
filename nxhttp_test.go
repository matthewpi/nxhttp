// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/matthewpi/nxhttp"
)

func ExampleFromClient() {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("Hello, world!"))
		}),
	)
	defer ts.Close()

	_ = nxhttp.FromClient(ts.Client())
}
