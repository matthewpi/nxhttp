// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxdial

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
)

// ErrInternalResolution is returned when a dialer attempts to connect to an
// internal IP address.
var ErrInternalResolution = errors.New("nxdial: destination resolves to an internal network location")

// RestrictedDialer is a [net.Dialer] wrapper that restricts the IP addresses
// that are allowed to be connected to.
//
// It is designed to be used with untrusted (user-provided) hostnames and avoid
// unintended internal access or information disclosure.
//
// It is safe from DNS rebinding attacks, hence why it is implemented at the
// dialer level. You are able to allow following HTTP redirects while using
// this dialer without compromising security.
//
// Example use-cases include but are not limited to: remote file downloads,
// calling user-provided webhook URLs, etc.
type RestrictedDialer struct {
	// dialer to forward calls to.
	dialer net.Dialer

	// AllowedPrefixes is a list of allowed [netip.Prefix].
	//
	// Any prefix present in the slice will be explicitly allowed no matter
	// what other options on the [RestrictedDialer] are configured, including
	// [RestrictedDialer.BlockedPrefixes].
	AllowedPrefixes []netip.Prefix

	// BlockedPrefixes is a list of blocked [netip.Prefix].
	//
	// Any prefix present here will be blocked unless there is an overlapping
	// prefix in [RestrictedDialer.AllowedPrefixes] in which case the
	// AllowedPrefixes option takes precedence.
	BlockedPrefixes []netip.Prefix

	// IsPrivate if enabled, blocks addresses that are considered private
	// according to [RFC 1918] for IPv4 addresses and [RFC 4193] for IPv6
	// addresses.
	//
	// Included prefixes:
	//
	// - 10.0.0.0/8 ([RFC 1918])
	// - 172.16.0.0/12 ([RFC 1918])
	// - 192.168.0.0/16 ([RFC 1918])
	// - fc00::/7 ([RFC 4193])
	//
	// [RFC 1918]: https://datatracker.ietf.org/doc/html/rfc1918
	// [RFC 4193]: https://datatracker.ietf.org/doc/html/rfc4193
	IsPrivate bool

	// IsLoopback if enabled, blocks loopback addresses.
	//
	// Included prefixes:
	//
	// - 127.0.0.0/8 ([RFC 1122])
	// - ::1/128 ([RFC 4291])
	//
	// [RFC 1122]: https://datatracker.ietf.org/doc/html/rfc1122#section-3.2.1.3
	// [RFC 4291]: https://datatracker.ietf.org/doc/html/rfc4291#section-2.4
	IsLoopback bool

	// IsLinkLocalUnicast if enabled, blocks link-local unicast addresses.
	//
	// Included prefixes:
	//
	// - 169.254.0.0/16 ([RFC 3927])
	// - fe80::/10 ([RFC 4291])
	//
	// [RFC 3927]: https://datatracker.ietf.org/doc/html/rfc3927#section-2.1
	// [RFC 4291]: https://datatracker.ietf.org/doc/html/rfc4291#section-2.4
	IsLinkLocalUnicast bool

	// IsLinkLocalMulticast if enabled, blocks link-local multicast addresses.
	IsLinkLocalMulticast bool

	// IsInterfaceLocalMulticast if enabled, blocks IPv6 interface-local
	// multicast addresses.
	IsInterfaceLocalMulticast bool
}

// NewRestrictedDialer returns a new [RestrictedDialer] with all predefined
// restrictions enabled.
//
// Callers are allowed to modify the returned [RestrictedDialer] before use
// to override the defaults or use other available options.
func NewRestrictedDialer() *RestrictedDialer {
	return &RestrictedDialer{
		IsPrivate:                 true,
		IsLoopback:                true,
		IsLinkLocalUnicast:        true,
		IsLinkLocalMulticast:      true,
		IsInterfaceLocalMulticast: true,
	}
}

// Dial connects to the address on the named network.
//
// See [net.Dial] for a description of the network and address
// parameters.
//
// Dial uses [context.Background] internally; to specify the context, use
// [Dialer.DialContext].
func (r *RestrictedDialer) Dial(network, addr string) (net.Conn, error) {
	return r.DialContext(context.Background(), network, addr)
}

// DialContext connects to the address on the named network using
// the provided context.
//
// The provided Context must be non-nil. If the context expires before
// the connection is complete, an error is returned. Once successfully
// connected, any expiration of the context will not affect the
// connection.
//
// When using TCP, and the host in the address parameter resolves to multiple
// network addresses, any dial timeout (from d.Timeout or ctx) is spread
// over each consecutive dial, such that each is given an appropriate
// fraction of the time to connect.
// For example, if a host has 4 IP addresses and the timeout is 1 minute,
// the connect to each single address will be given 15 seconds to complete
// before trying the next one.
//
// See [net.Dial] for a description of the network and address parameters.
func (r *RestrictedDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	// Forward the connection to the underlying dialer.
	c, err := r.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	// Parse the IP address and port we are connecting to.
	addrPort, err := netip.ParseAddrPort(c.RemoteAddr().String())
	if err != nil {
		return c, fmt.Errorf("nxhttp: failed to parse remote address: %w", err)
	}

	// Check if the address is restricted.
	if !r.IsAllowed(addrPort.Addr()) {
		return c, ErrInternalResolution
	}

	return c, nil
}

// IsAllowed checks if addr is allowed to be dialed as per the restrictions
// of the dialer.
//
// Returns `true` if addr is allowed, `false` otherwise.
func (r *RestrictedDialer) IsAllowed(addr netip.Addr) bool {
	// If the address is within one of the allowed prefixes, allow it and skip
	// any further checks.
	for _, p := range r.AllowedPrefixes {
		if p.Contains(addr) {
			return true
		}
	}

	// If the address is within one of the blocked blocks, deny it and skip
	// any further checks.
	for _, p := range r.BlockedPrefixes {
		if p.Contains(addr) {
			return false
		}
	}

	if r.IsPrivate && addr.IsPrivate() {
		return false
	}

	if r.IsLoopback && addr.IsLoopback() {
		return false
	}

	if r.IsLinkLocalUnicast && addr.IsLinkLocalUnicast() {
		return false
	}

	if r.IsLinkLocalMulticast && addr.IsLinkLocalMulticast() {
		return false
	}

	if r.IsInterfaceLocalMulticast && addr.IsInterfaceLocalMulticast() {
		return false
	}

	// The address is allowed.
	return true
}
