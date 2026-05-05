package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStartLocalServerUsesRepoRootedContentDir(t *testing.T) {
	root, err := repoRoot()
	if err != nil {
		t.Fatal(err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	if err := os.Chdir(filepath.Join(root, "cmd", "visual-smoke")); err != nil {
		t.Fatal(err)
	}

	baseURL, cleanup, err := startLocalServer(root)
	if err != nil {
		t.Fatalf("start local server: %v", err)
	}
	t.Cleanup(cleanup)

	res, err := http.Get(baseURL + "/docs")
	if err != nil {
		t.Fatalf("get docs landing: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read docs landing: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("/docs status = %d, want %d: %s", res.StatusCode, http.StatusOK, string(body))
	}
	if !strings.Contains(string(body), "File-owned landing content") {
		t.Fatalf("/docs did not render file-backed content: %s", string(body))
	}
}
