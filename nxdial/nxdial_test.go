// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: Copyright (c) 2026 Matthew Penner

package nxdial_test

import (
	"net/netip"
	"testing"

	"github.com/matthewpi/nxhttp/nxdial"
)

func TestRestrictedDialer(t *testing.T) {
	t.Run("AllowedPrefixes and BlockedPrefixes", func(t *testing.T) {
		d := &nxdial.RestrictedDialer{
			AllowedPrefixes: []netip.Prefix{
				netip.MustParsePrefix("1.1.1.1/32"),
				netip.MustParsePrefix("2606:4700:4700::1111/128"),
			},
			BlockedPrefixes: []netip.Prefix{
				netip.MustParsePrefix("0.0.0.0/0"),
				netip.MustParsePrefix("::/0"),
			},
		}
		for i, tc := range []struct {
			addr string
			ok   bool
		}{
			// Ensure our addresses are allowed.
			{"1.1.1.1", true},
			{"2606:4700:4700::1111", true},

			// Ensure other addresses are blocked.
			{"127.0.0.1", false},
			{"::1", false},
			{"1.0.0.1", false},
			{"2606:4700:4700::1001", false},
		} {
			addr, err := netip.ParseAddr(tc.addr)
			if err != nil {
				t.Errorf("netip.ParseAddr(%q) #%d: %v", tc.addr, i, err)
				return
			}

			if d.IsAllowed(addr) != tc.ok {
				t.Errorf("IsAllowed(%q) #%d: expected %t, but got %t", tc.addr, i, tc.ok, !tc.ok)
			}
		}
	})

	t.Run("IsPrivate", func(t *testing.T) {
		d := &nxdial.RestrictedDialer{IsPrivate: true}
		for i, tc := range []struct {
			addr string
			ok   bool
		}{
			// Ensure private addresses are blocked.
			{"10.0.0.1", false},
			{"172.16.0.1", false},
			{"192.168.0.1", false},
			{"fc00::1", false},

			// Ensure public addresses are still allowed.
			{"1.1.1.1", true},
			{"2606:4700:4700::1111", true},

			// Ensure loopback addresses are still allowed.
			{"127.0.0.1", true},
			{"127.0.0.2", true},
			{"::1", true},

			// Ensure link-local unicast addresses are still allowed.
			{"169.254.0.1", true},
			{"fe80::1", true},
		} {
			addr, err := netip.ParseAddr(tc.addr)
			if err != nil {
				t.Errorf("netip.ParseAddr(%q) #%d: %v", tc.addr, i, err)
				return
			}

			if d.IsAllowed(addr) != tc.ok {
				t.Errorf("IsAllowed(%q) #%d: expected %t, but got %t", tc.addr, i, tc.ok, !tc.ok)
			}
		}
	})

	t.Run("IsLoopback", func(t *testing.T) {
		d := &nxdial.RestrictedDialer{IsLoopback: true}
		for i, tc := range []struct {
			addr string
			ok   bool
		}{
			// Ensure loopback addresses are blocked.
			{"127.0.0.1", false},
			{"127.0.0.2", false},
			{"::1", false},

			// Ensure public addresses are still allowed.
			{"1.1.1.1", true},
			{"2606:4700:4700::1111", true},

			// Ensure private addresses are still allowed.
			{"10.0.0.1", true},
			{"172.16.0.1", true},
			{"192.168.0.1", true},
			{"fc00::1", true},

			// Ensure link-local unicast addresses are still allowed.
			{"169.254.0.1", true},
			{"fe80::1", true},
		} {
			addr, err := netip.ParseAddr(tc.addr)
			if err != nil {
				t.Errorf("netip.ParseAddr(%q) #%d: %v", tc.addr, i, err)
				return
			}

			if d.IsAllowed(addr) != tc.ok {
				t.Errorf("IsAllowed(%q) #%d: expected %t, but got %t", tc.addr, i, tc.ok, !tc.ok)
			}
		}
	})

	t.Run("IsLinkLocalUnicast", func(t *testing.T) {
		d := &nxdial.RestrictedDialer{IsLinkLocalUnicast: true}
		for i, tc := range []struct {
			addr string
			ok   bool
		}{
			// Ensure link-local unicast addresses are blocked.
			{"169.254.0.1", false},
			{"fe80::1", false},

			// Ensure public addresses are still allowed.
			{"1.1.1.1", true},
			{"2606:4700:4700::1111", true},

			// Ensure loopback addresses are still allowed.
			{"127.0.0.1", true},
			{"127.0.0.2", true},
			{"::1", true},

			// Ensure private addresses are still allowed.
			{"10.0.0.1", true},
			{"172.16.0.1", true},
			{"192.168.0.1", true},
			{"fc00::1", true},
		} {
			addr, err := netip.ParseAddr(tc.addr)
			if err != nil {
				t.Errorf("netip.ParseAddr(%q) #%d: %v", tc.addr, i, err)
				return
			}

			if d.IsAllowed(addr) != tc.ok {
				t.Errorf("IsAllowed(%q) #%d: expected %t, but got %t", tc.addr, i, tc.ok, !tc.ok)
			}
		}
	})

	// TODO: IsLinkLocalMulticast

	// TODO: IsInterfaceLocalMulticast
}
