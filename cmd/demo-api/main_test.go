package main

import (
	"net"
	"testing"

	"sf-mvp/internal/httpapi"
)

func TestNewServerUsesDefaultLoopbackAddress(t *testing.T) {
	server, err := newServer("")
	if err != nil {
		t.Fatalf("newServer returned error: %v", err)
	}

	if server.Addr != httpapi.DefaultListenAddress() {
		t.Fatalf("server.Addr = %q, want %q", server.Addr, httpapi.DefaultListenAddress())
	}
	host, _, err := net.SplitHostPort(server.Addr)
	if err != nil {
		t.Fatalf("server.Addr is not host:port: %v", err)
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		t.Fatalf("server.Addr host = %q, want loopback IP", host)
	}
	if server.Handler == nil {
		t.Fatal("server.Handler is nil")
	}
}

func TestNewServerAllowsLoopbackAddressOverride(t *testing.T) {
	server, err := newServer("127.0.0.1:18080")
	if err != nil {
		t.Fatalf("newServer returned error: %v", err)
	}

	if server.Addr != "127.0.0.1:18080" {
		t.Fatalf("server.Addr = %q, want 127.0.0.1:18080", server.Addr)
	}
}

func TestNewServerRejectsNonLoopbackAddressOverride(t *testing.T) {
	_, err := newServer("0.0.0.0:8080")
	if err == nil {
		t.Fatal("newServer returned nil error, want non-loopback rejection")
	}
}
