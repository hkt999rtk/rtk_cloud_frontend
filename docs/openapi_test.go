package docs_test

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPIDocumentCoversRuntimeEndpoints(t *testing.T) {
	contents, err := os.ReadFile("openapi.yaml")
	if err != nil {
		t.Fatalf("read openapi.yaml: %v", err)
	}

	var document struct {
		OpenAPI string                    `yaml:"openapi"`
		Info    map[string]any            `yaml:"info"`
		Paths   map[string]map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(contents, &document); err != nil {
		t.Fatalf("parse openapi.yaml: %v", err)
	}

	if document.OpenAPI != "3.1.0" {
		t.Fatalf("openapi = %q, want 3.1.0", document.OpenAPI)
	}
	if document.Info["title"] == "" {
		t.Fatal("info.title is required")
	}

	expected := map[string][]string{
		"/api/event":            {"post"},
		"/api/search":           {"post"},
		"/contact":              {"get", "post"},
		"/zh-tw/contact":        {"get", "post"},
		"/zh-cn/contact":        {"get", "post"},
		"/healthz":              {"get"},
		"/admin/leads":          {"get"},
		"/admin/leads.csv":      {"get"},
		"/admin/reload-content": {"post"},
	}
	for path, methods := range expected {
		operations, ok := document.Paths[path]
		if !ok {
			t.Fatalf("missing path %s", path)
		}
		for _, method := range methods {
			if operations[method] == nil {
				t.Fatalf("missing %s %s", method, path)
			}
		}
	}
}
