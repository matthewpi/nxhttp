// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

// Package httpheader provides constantly typed HTTP headers in their
// canonicalized form to avoid additional string allocations and the
// possibility of typos.
package httpheader

import (
	"net/http"
	"net/textproto"
)

// Key represents the canonical format of a MIME header key.
type Key string

// Canonicalize returns the canonical format of the MIME header key s.
//
// This is the same as [textproto.CanonicalMIMEHeaderKey] except the return
// type is a [Key] and not a [string].
func Canonicalize(s string) Key {
	return Key(textproto.CanonicalMIMEHeaderKey(s))
}

// Add is like [http.Header.Add], but the key must already be in [Key] form.
func Add(h http.Header, key Key, value string) {
	if h == nil {
		return
	}
	h[string(key)] = append(h[string(key)], value)
}

// Set is like [http.Header.Set], but the key must already be in [Key] form.
func Set(h http.Header, key Key, value ...string) {
	if h == nil {
		return
	}
	h[string(key)] = value
}

// Get is like [http.Header.Get], but the key must already be in [Key] form.
func Get(h http.Header, key Key) string {
	if h == nil {
		return ""
	}
	v := h[string(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Values is like [http.Header.Values], but the key must already be in [Key] form.
func Values(h http.Header, key Key) []string {
	if h == nil {
		return nil
	}
	return h[string(key)]
}

// Del is like [http.Header.Del], but the key must already be in [Key] form.
func Del(h http.Header, key Key) {
	if h == nil {
		return
	}
	delete(h, string(key))
}

const (
	// Accept is the HTTP [Accept] header.
	// [Accept]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept
	Accept Key = "Accept"

	// AcceptEncoding is the HTTP [Accept-Encoding] header.
	// [Accept-Encoding]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Encoding
	AcceptEncoding Key = "Accept-Encoding"

	// AcceptRanges is the HTTP [Accept-Ranges] header.
	// [Accept-Ranges]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Accept-Ranges
	AcceptRanges Key = "Accept-Ranges"

	// Age is the HTTP [Age] header.
	// [Age]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Age
	Age Key = "Age"

	// Allow is the HTTP [Allow] header.
	// [Allow]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Allow
	Allow Key = "Allow"

	// AltSvc is the HTTP [Alt-Svc] header.
	// [Alt-Svc]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Alt-Svc
	AltSvc Key = "Alt-Svc"

	// AltUsed is the HTTP [Alt-Used] header.
	// [Alt-Used]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Alt-Used
	AltUsed Key = "Alt-Used"

	// Authorization is the HTTP [Authorization] header.
	// [Authorization]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Authorization
	Authorization Key = "Authorization"

	// CacheControl is the HTTP [Cache-Control] header.
	// [Cache-Control]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cache-Control
	CacheControl Key = "Cache-Control"

	// Connection is the HTTP [Connection] header.
	// [Connection]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Connection
	Connection Key = "Connection"

	// ContentDigest is the HTTP [Content-Digest] header.
	// [Content-Digest]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Digest
	ContentDigest Key = "Content-Digest"

	// ContentDisposition is the HTTP [Content-Disposition] header.
	// [Content-Disposition]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Disposition
	ContentDisposition Key = "Content-Disposition"

	// ContentEncoding is the HTTP [Content-Encoding] header.
	// [Content-Encoding]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Encoding
	ContentEncoding Key = "Content-Encoding"

	// ContentLanguage is the HTTP [Content-Language] header.
	// [Content-Language]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Language
	ContentLanguage Key = "Content-Language"

	// ContentLength is the HTTP [Content-Length] header.
	// [Content-Length]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Length
	ContentLength Key = "Content-Length"

	// ContentLocation is the HTTP [Content-Location] header.
	// [Content-Location]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Location
	ContentLocation Key = "Content-Location"

	// ContentRange is the HTTP [Content-Range] header.
	// [Content-Range]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Range
	ContentRange Key = "Content-Range"

	// ContentSecurityPolicy is the HTTP [Content-Security-Policy] header.
	// [Content-Security-Policy]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Security-Policy
	ContentSecurityPolicy Key = "Content-Security-Policy"

	// ContentType is the HTTP [Content-Type] header.
	// [Content-Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Type
	ContentType Key = "Content-Type"

	// Cookie is the HTTP [Cookie] header.
	// [Cookie]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Cookie
	Cookie Key = "Cookie"

	// Date is the HTTP [Date] header.
	// [Date]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Date
	Date Key = "Date"

	// EarlyData is the HTTP [Early-Data] header.
	// [Early-Data]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Early-Data
	EarlyData Key = "Early-Data"

	// ETag is the HTTP [ETag] header.
	// [ETag]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/ETag
	ETag Key = "Etag" // lower-case t is intended, do not change it.

	// Expect is the HTTP [Expect] header.
	// [Expect]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expect
	Expect Key = "Expect"

	// Expires is the HTTP [Expires] header.
	// [Expires]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Expires
	Expires Key = "Expires"

	// Forwarded is the HTTP [Forwarded] header.
	// [Forwarded]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Forwarded
	Forwarded Key = "Forwarded"

	// From is the HTTP [From] header.
	// [From]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/From
	From Key = "From"

	// Host is the HTTP [Host] header.
	// [Host]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Host
	Host Key = "Host"

	// IdempotencyKey is the HTTP [Idempotency-Key] header.
	// [Idempotency-Key]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Idempotency-Key
	IdempotencyKey Key = "Idempotency-Key"

	// IfMatch is the HTTP [If-Match] header.
	// [If-Match]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Match
	IfMatch Key = "If-Match"

	// IfModifiedSince is the HTTP [If-Modified-Since] header.
	// [If-Modified-Since]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Modified-Since
	IfModifiedSince Key = "If-Modified-Since"

	// IfNoneMatch is the HTTP [If-None-Match] header.
	// [If-None-Match]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-None-Match
	IfNoneMatch Key = "If-None-Match"

	// IfRange is the HTTP [If-Range] header.
	// [If-Range]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Range
	IfRange Key = "If-Range"

	// IfUnmodifiedSince is the HTTP [If-Unmodified-Since] header.
	// [If-Unmodified-Since]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/If-Unmodified-Since
	IfUnmodifiedSince Key = "If-Unmodified-Since"

	// KeepAlive is the HTTP [Keep-Alive] header.
	// [Keep-Alive]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Keep-Alive
	KeepAlive Key = "Keep-Alive"

	// LastModified is the HTTP [Last-Modified] header.
	// [Last-Modified]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Last-Modified
	LastModified Key = "Last-Modified"

	// Link is the HTTP [Link] header.
	// [Link]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Link
	Link Key = "Link"

	// Location is the HTTP [Location] header.
	// [Location]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Location
	Location Key = "Location"

	// Origin is the HTTP [Origin] header.
	// [Origin]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Origin
	Origin Key = "Origin"

	// Range is the HTTP [Range] header.
	// [Range]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Range
	Range Key = "Range"

	// ReprDigest is the HTTP [Repr-Digest] header.
	// [Repr-Digest]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Repr-Digest
	ReprDigest Key = "Repr-Digest"

	// RetryAfter is the HTTP [Retry-After] header.
	// [Retry-After]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Retry-After
	RetryAfter Key = "Retry-After"

	// Server is the HTTP [Server] header.
	// [Server]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Server
	Server Key = "Server"

	// SetCookie is the HTTP [Set-Cookie] header.
	// [Set-Cookie]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Set-Cookie
	SetCookie Key = "Set-Cookie"

	// TE is the HTTP [TE] header.
	// [TE]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/TE
	TE Key = "Te" // lower-case e is intended, do not change it.

	// Trailer is the HTTP [Trailer] header.
	// [Trailer]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Trailer
	Trailer Key = "Trailer"

	// TransferEncoding is the HTTP [Transfer-Encoding] header.
	// [Transfer-Encoding]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Transfer-Encoding
	TransferEncoding Key = "Transfer-Encoding"

	// Upgrade is the HTTP [Upgrade] header.
	// [Upgrade]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Upgrade
	Upgrade Key = "Upgrade"

	// UserAgent is the HTTP [User-Agent] header.
	// [User-Agent]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/User-Agent
	UserAgent Key = "User-Agent"

	// Vary is the HTTP [Vary] header.
	// [Vary]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Vary
	Vary Key = "Vary"

	// Via is the HTTP [Via] header.
	// [Via]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Via
	Via Key = "Via"

	// WWWAuthenticate is the HTTP [WWW-Authenticate] header.
	// [WWW-Authenticate]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/WWW-Authenticate
	WWWAuthenticate Key = "Www-Authenticate" // lower-case "w"s are intended, do not change it.
)

func init() {
	// In order to allow those keys to be constants they must be bare strings,
	// but we want to make 100% sure they are properly canonicalized to avoid
	// any problems.
	for _, got := range []Key{
		Accept,
		AcceptEncoding,
		AcceptRanges,
		Age,
		Allow,
		AltSvc,
		AltUsed,
		Authorization,
		CacheControl,
		Connection,
		ContentDigest,
		ContentDisposition,
		ContentEncoding,
		ContentLanguage,
		ContentLength,
		ContentLocation,
		ContentRange,
		ContentSecurityPolicy,
		ContentType,
		Cookie,
		Date,
		EarlyData,
		ETag,
		Expect,
		Expires,
		Forwarded,
		From,
		Host,
		IdempotencyKey,
		IfMatch,
		IfModifiedSince,
		IfNoneMatch,
		IfRange,
		IfUnmodifiedSince,
		KeepAlive,
		LastModified,
		Link,
		Location,
		Origin,
		Range,
		ReprDigest,
		RetryAfter,
		Server,
		SetCookie,
		TE,
		Trailer,
		TransferEncoding,
		Upgrade,
		UserAgent,
		Vary,
		Via,
		WWWAuthenticate,
	} {
		if expected := Canonicalize(string(got)); got != expected {
			panic("nxhttp/httpheader: header is not properly canonicalized (got: \"" + got + "\", expected: \"" + expected + "\")")
		}
	}
}
