// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package httpheader

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ParseRetryAfter parses an HTTP [Retry-After] header into a [time.Duration].
//
// The value of the header could be either an HTTP Date (see [http.ParseTime])
// or a number of seconds.
//
// [Retry-After]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After
func ParseRetryAfter(v string) (time.Duration, error) {
	// Fast-path, empty string.
	if v == "" {
		return 0, nil
	}

	// Most services set Retry-After to a number of seconds, so attempt to
	// parse it as an integer first.
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		if i < 0 {
			return 0, fmt.Errorf("nxhttp: got negative Retry-After value in response (%d)", i)
		}
		// The string is an integer, convert it to a [time.Duration] with a unit
		// of [time.Second].
		return time.Duration(i) * time.Second, nil
	}

	// Retry-After is either malformed or is a Date, so attempt to parse it.
	t, err := http.ParseTime(v)
	if err != nil {
		// Retry-After is neither the number of seconds or a date.
		return 0, fmt.Errorf("nxhttp: failed to parse Retry-After header from response: %w", err)
	}

	// We want to return a [time.Duration], not a [time.Time] to the caller, so
	// subtract the current time from the Retry-After time.
	until := time.Until(t)

	// Ensure the duration isn't negative as most callers would not be expecting
	// a negative value.
	if until < 0 {
		return 0, fmt.Errorf("nxhttp: Retry-After date is in the past (%s)", t)
	}

	// Return the calculated [time.Duration] from the date.
	return until, nil
}
