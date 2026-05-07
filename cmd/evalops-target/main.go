package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"sf-mvp/internal/eval"
)

const defaultListenAddr = "127.0.0.1:18085"

func main() {
	addr := flag.String("addr", defaultListenAddr, "loopback listen address")
	flag.Parse()

	server, err := newServer(*addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("evalops target listening on http://%s/evalops/incident", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func newServer(addr string) (*http.Server, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = defaultListenAddr
	}
	if err := validateLoopbackAddress(addr); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/evalops/incident", eval.NewIncidentEvalTarget())
	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}, nil
}

func validateLoopbackAddress(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("listen address must be host:port: %w", err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("listen address must use a loopback IP: %s", addr)
	}
	return nil
}
