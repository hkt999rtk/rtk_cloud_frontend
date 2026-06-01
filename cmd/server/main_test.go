package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
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
		done <- serveWithGracefulShutdown(ctx, server, zap.NewNop(), time.Second)
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
		done <- serveWithGracefulShutdown(ctx, server, zap.NewNop(), time.Second)
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
		"COPY --from=builder /src/content /app/content",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("Dockerfile missing %q", expected)
		}
	}
}

func TestCDWorkflowPackagesDataForServiceUser(t *testing.T) {
	contents, err := os.ReadFile("../../.github/workflows/cd.yml")
	if err != nil {
		t.Fatalf("read cd workflow: %v", err)
	}

	text := string(contents)
	for _, expected := range []string{
		`service_user="$(systemctl show "$service_name" --property=User --value)"`,
		`service_group="$(systemctl show "$service_name" --property=Group --value)"`,
		`--owner="$service_user"`,
		`--group="$service_group"`,
		`--mode=u+rwX,go+rX`,
		`cp -R content templates static dist/`,
		`test -f dist/content/docs/en/docs.yaml`,
		`tar "${tar_owner_args[@]}" -C dist -czf - bin content templates static data | sudo /usr/local/sbin/deploy-realtek-connect`,
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("cd workflow missing %q", expected)
		}
	}
}

func TestInstallScriptDefinesRuntimeLogForwarderLabels(t *testing.T) {
	contents, err := os.ReadFile("../../deploy/install.sh")
	if err != nil {
		t.Fatalf("read install script: %v", err)
	}

	text := string(contents)
	for _, expected := range []string{
		"Environment=REALTEK_CONNECT_VERSION=$version",
		"Environment=RTK_LOG_FORWARDER_JOURNAL_LABELS=service=realtek-connect,unit=realtek-connect.service,component=server",
		"Environment=RTK_LOG_FORWARDER_NGINX_ACCESS_LABELS=service=realtek-connect,unit=nginx.service,component=nginx-access",
		"Environment=RTK_LOG_FORWARDER_NGINX_ERROR_LABELS=service=realtek-connect,unit=nginx.service,component=nginx-error",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("install script missing %q", expected)
		}
	}
}
