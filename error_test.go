// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/matthewpi/nxhttp"
)

func TestStatusError(t *testing.T) {
	var err error
	err = nxhttp.NewStatusError(nil, http.StatusOK)
	if sErr, ok := nxhttp.AsStatusError(err); !ok {
		t.Error("return result of nxhttp.NewStatusError does not work with nxhttp.AsStatusError")
	} else if sErr.Expected != http.StatusOK {
		t.Error("nxhttp.AsStatusError returned a different nxhttp.StatusError")
	}

	// Wrap the error so we can ensure [nxhttp.AsStatusError] still works even
	// if the error is wrapped.
	err = fmt.Errorf("%w", err)

	if sErr, ok := nxhttp.AsStatusError(err); !ok {
		t.Error("wrapped result of nxhttp.NewStatusError does not work with nxhttp.AsStatusError")
	} else if sErr.Expected != http.StatusOK {
		t.Error("nxhttp.AsStatusError returned a different nxhttp.StatusError")
	}
}
