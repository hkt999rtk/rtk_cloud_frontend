package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type fakeServer struct {
	listenStarted  chan struct{}
	shutdownCalled chan struct{}
	listenErr      error
	shutdownErr    error
}

func newFakeServer(listenErr error, shutdownErr error) *fakeServer {
	return &fakeServer{
		listenStarted:  make(chan struct{}),
		shutdownCalled: make(chan struct{}),
		listenErr:      listenErr,
		shutdownErr:    shutdownErr,
	}
}

func (s *fakeServer) ListenAndServe() error {
	close(s.listenStarted)
	<-s.shutdownCalled
	if s.listenErr != nil {
		return s.listenErr
	}
	return http.ErrServerClosed
}

func (s *fakeServer) Shutdown(context.Context) error {
	close(s.shutdownCalled)
	return s.shutdownErr
}

func TestServeWithGracefulShutdownStopsServerOnContextCancel(t *testing.T) {
	server := newFakeServer(nil, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- serveWithGracefulShutdown(ctx, server, log.New(io.Discard, "", 0), time.Second)
	}()

	<-server.listenStarted
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serveWithGracefulShutdown returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serveWithGracefulShutdown did not return")
	}
}

func TestServeWithGracefulShutdownReturnsListenErrors(t *testing.T) {
	server := newFakeServer(errors.New("listen failed"), nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- serveWithGracefulShutdown(ctx, server, log.New(io.Discard, "", 0), time.Second)
	}()

	<-server.listenStarted
	close(server.shutdownCalled)

	select {
	case err := <-done:
		if err == nil || err.Error() != "listen failed" {
			t.Fatalf("error = %v, want listen failed", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serveWithGracefulShutdown did not return")
	}
}

func TestDockerfileUsesPersistentSQLitePaths(t *testing.T) {
	contents, err := os.ReadFile("../../Dockerfile")
	if err != nil {
		t.Fatalf("read Dockerfile: %v", err)
	}

	text := string(contents)
	for _, expected := range []string{
		"ENV DATABASE_PATH=/data/connectplus.db",
		"ENV ANALYTICS_DATABASE_PATH=/data/analytics.db",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("Dockerfile missing %q", expected)
		}
	}
}
