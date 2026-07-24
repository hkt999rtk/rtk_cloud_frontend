package search

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAIClientEmbedAndAnswer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" || r.Method != http.MethodPost {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/v1/embeddings":
			_, _ = w.Write([]byte(`{"data":[{"index":1,"embedding":[0,1]},{"index":0,"embedding":[1,0]}]}`))
		case "/v1/responses":
			_, _ = w.Write([]byte(`{"output_text":"documented answer"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := OpenAIClient{APIKey: "test-key", BaseURL: server.URL + "/", HTTPClient: server.Client()}

	vectors, err := client.Embed(context.Background(), []string{"one", "two"})
	if err != nil {
		t.Fatal(err)
	}
	if len(vectors) != 2 || vectors[0][0] != 1 || vectors[1][1] != 1 {
		t.Fatalf("vectors = %#v", vectors)
	}
	answer, err := client.Answer(context.Background(), AnswerRequest{Query: "What?", Locale: "en", Context: "source"})
	if err != nil || answer != "documented answer" {
		t.Fatalf("Answer() = %q, %v", answer, err)
	}
}

func TestOpenAIClientAnswerContentFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"output":[{"content":[{"type":"output_text","text":"fallback answer"}]}]}`))
	}))
	defer server.Close()
	client := OpenAIClient{APIKey: "key", BaseURL: server.URL}
	answer, err := client.Answer(context.Background(), AnswerRequest{})
	if err != nil || answer != "fallback answer" {
		t.Fatalf("Answer() = %q, %v", answer, err)
	}
}

func TestOpenAIClientErrorFamilies(t *testing.T) {
	if _, err := (OpenAIClient{}).Embed(context.Background(), []string{"one"}); err == nil {
		t.Fatal("Embed accepted an empty API key")
	}
	if _, err := (OpenAIClient{}).Answer(context.Background(), AnswerRequest{}); err == nil {
		t.Fatal("Answer accepted an empty API key")
	}
	tests := []struct {
		name string
		body string
		code int
		run  func(OpenAIClient) error
	}{
		{"missing vector", `{"data":[]}`, http.StatusOK, func(c OpenAIClient) error { _, err := c.Embed(context.Background(), []string{"one"}); return err }},
		{"missing answer", `{"output":[]}`, http.StatusOK, func(c OpenAIClient) error { _, err := c.Answer(context.Background(), AnswerRequest{}); return err }},
		{"invalid json", `{`, http.StatusOK, func(c OpenAIClient) error { _, err := c.Embed(context.Background(), []string{"one"}); return err }},
		{"api status", `provider credential must not be copied`, http.StatusBadGateway, func(c OpenAIClient) error { _, err := c.Embed(context.Background(), []string{"one"}); return err }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.code)
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()
			err := test.run(OpenAIClient{APIKey: "key", BaseURL: server.URL})
			if err == nil {
				t.Fatal("request unexpectedly passed")
			}
		})
	}
	transportErr := errors.New("offline")
	client := OpenAIClient{
		APIKey: "key",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, transportErr
		})},
	}
	if _, err := client.Embed(context.Background(), []string{"one"}); !errors.Is(err, transportErr) {
		t.Fatalf("transport error = %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestBuildIndexChunksValidation(t *testing.T) {
	embedder := &stubEmbedder{vectors: map[string][]float64{}}
	chunks, err := BuildIndexChunks(context.Background(), []Document{
		{ID: "", Body: "ignored"},
		{ID: "doc", Title: "Document", Body: strings.Repeat("content ", 300)},
	}, embedder)
	if err != nil || len(chunks) < 2 || embedder.calls != len(chunks) {
		t.Fatalf("BuildIndexChunks() chunks=%d calls=%d err=%v", len(chunks), embedder.calls, err)
	}
	if chunks[0].ChunkID != "doc:1" {
		t.Fatalf("first chunk = %#v", chunks[0])
	}
	empty, err := BuildIndexChunks(context.Background(), nil, embedder)
	if err != nil || empty != nil {
		t.Fatalf("empty chunks = %#v, %v", empty, err)
	}
	bad := mismatchEmbedder{}
	if _, err := BuildIndexChunks(context.Background(), []Document{{ID: "doc", Body: "body"}}, bad); err == nil {
		t.Fatal("embedding count mismatch passed")
	}
}

type mismatchEmbedder struct{}

func (mismatchEmbedder) Embed(context.Context, []string) ([][]float64, error) {
	return nil, nil
}
