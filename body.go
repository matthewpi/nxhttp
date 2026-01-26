// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strings"
)

// BodyFunc is a function that can be used with [http.Request.GetBody].
//
// The [net/http] package uses it for allowing request bodies to be sent even
// through redirects, however we use it to allow for retrying failed requests.
//
// [net/http]: https://pkg.go.dev/net/http
type BodyFunc func() (io.ReadCloser, error)

// ReadOpener is an interface for anything that can open a new [io.ReadCloser].
type ReadOpener interface {
	// Open opens a new reader.
	Open() (io.ReadCloser, error)
}

// readOpener is a simple implementation of the [ReadOpener] interface.
type readOpener struct {
	fn BodyFunc
	n  int64
}

func (r *readOpener) Open() (io.ReadCloser, error) { return r.fn() }
func (r *readOpener) Size() int64                  { return r.n }

// ReadOpenerFor creates a new [ReadOpener] that uses fn as the opener and n
// for the size. If the size of the data is unknown, use `-1` as the value
// for n. Only use `0` as the value for n if the length of the data from the
// reader is actually 0.
func ReadOpenerFor(fn BodyFunc, n int64) ReadOpener {
	return &readOpener{
		fn: fn,
		n:  n,
	}
}

// SeekableFile is an interface for an [fs.File] that also implements
// [io.Seeker].
type SeekableFile interface {
	fs.File
	io.Seeker
}

// GetBody returns a [BodyFunc] and size for the given body `v`. If `v` is not
// a recognized type, an error will be returned.
func GetBody(v any) (BodyFunc, int64, error) {
	if v == nil {
		return nil, 0, nil
	}
	switch body := v.(type) {
	case ReadOpener:
		return body.Open, getLen(body), nil
	case SeekableFile: // [fs.File] with [io.Seeker]
		s, err := body.Stat()
		if err != nil {
			return nil, 0, fmt.Errorf("nxhttp: failed to stat file to use as body: %w", err)
		}
		return bodyFuncFromFile(body, s.Size())
	case *os.File:
		s, err := body.Stat()
		if err != nil {
			return nil, 0, fmt.Errorf("nxhttp: failed to stat file to use as body: %w", err)
		}
		return bodyFuncFromFile(body, s.Size())
	case io.ReadSeeker:
		return bodyFuncFromReadSeeker(body), getLen(body), nil
	case []byte:
		return bodyFuncFromReadSeekerSize(bytes.NewReader(body))
	case string:
		return bodyFuncFromReadSeekerSize(strings.NewReader(body))
	case url.Values:
		return bodyFuncFromReadSeekerSize(strings.NewReader(body.Encode()))
	default:
		return nil, 0, fmt.Errorf("nxhttp: cannot handle body of type %T", v)
	}
}

// bodyFuncFromFile returns a reusable [BodyFunc] given an...
func bodyFuncFromFile(r io.ReadSeeker, n int64) (BodyFunc, int64, error) {
	return func() (io.ReadCloser, error) {
		// Seek to the beginning.
		_, err := r.Seek(0, io.SeekStart)

		// Since files can be written to while we are reading, limit our
		// read to the size of the file we got. This avoids us writing more
		// data than we reported as the ContentLength.
		return io.NopCloser(io.LimitReader(r, n)), err
	}, n, nil
}

// bodyFuncFromReadSeeker returns a reusable [BodyFunc] given an
// [io.ReadSeeker].
func bodyFuncFromReadSeeker(r io.ReadSeeker) BodyFunc {
	return func() (io.ReadCloser, error) {
		// Seek to the beginning.
		_, err := r.Seek(0, io.SeekStart)
		return io.NopCloser(r), err
	}
}

type readSeekerSize interface {
	io.ReadSeeker
	Size() int64
}

// bodyFuncFromReadSeekerSize returns a reusable [BodyFunc] and size given a
// [ReadSeekerSize].
func bodyFuncFromReadSeekerSize(r readSeekerSize) (BodyFunc, int64, error) {
	return bodyFuncFromReadSeeker(r), r.Size(), nil
}

// getLen attempts to get the length from a struct by performing interface
// assertions for different `Len` and `Size` methods.
func getLen(v any) int64 {
	switch r := v.(type) {
	case interface{ Size() int64 }:
		// This covers a few types.
		//
		// - [*bytes.Reader]
		// - [*io.SectionReader]
		// - [*strings.Reader]
		//
		// NOTE: technically we would want to prefer Len in order to properly
		// handle a partially read buffer. But if we are given an [io.ReadSeeker]
		// like the types listed above are, we will seek to the start of the
		// buffer on each attempt, so using Len would be inaccurate.
		return r.Size()
	case interface{ Len() int64 }:
		return r.Len()
	case interface{ Len() int }:
		return int64(r.Len())
	default:
		return -1
	}
}
