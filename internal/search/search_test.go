package search

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

type stubEmbedder struct {
	vectors map[string][]float64
	calls   int
}

func (s *stubEmbedder) Embed(ctx context.Context, inputs []string) ([][]float64, error) {
	s.calls += len(inputs)
	out := make([][]float64, 0, len(inputs))
	for _, input := range inputs {
		if vector, ok := s.vectors[input]; ok {
			out = append(out, vector)
			continue
		}
		out = append(out, []float64{0, 1, 0})
	}
	return out, nil
}

type stubAnswerer struct {
	calls   int
	context string
}

func (s *stubAnswerer) Answer(ctx context.Context, req AnswerRequest) (string, error) {
	s.calls++
	s.context = req.Context
	return "OTA supports firmware rollout based on the cited documentation.", nil
}

func TestRepositoryStoresAndSearchesEmbeddings(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()

	docs := []Document{
		{ID: "feature:ota:en", Locale: "en", SourceType: "feature", Title: "OTA", URL: "/features/ota", Body: "OTA firmware rollout campaign controls."},
		{ID: "feature:provision:en", Locale: "en", SourceType: "feature", Title: "Provision", URL: "/features/provision", Body: "Device onboarding and activation."},
	}
	if err := repo.Replace(ctx, []IndexedChunk{
		{Document: docs[0], ChunkID: "ota-1", Text: docs[0].Body, Embedding: []float64{1, 0, 0}},
		{Document: docs[1], ChunkID: "provision-1", Text: docs[1].Body, Embedding: []float64{0, 1, 0}},
	}); err != nil {
		t.Fatalf("replace index: %v", err)
	}

	results, err := repo.Search(ctx, []float64{1, 0, 0}, 3)
	if err != nil {
		t.Fatalf("search index: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected search results")
	}
	if results[0].Title != "OTA" || results[0].URL != "/features/ota" {
		t.Fatalf("unexpected top result: %+v", results[0])
	}
	if results[0].Score < 0.99 {
		t.Fatalf("score = %f, want high cosine score", results[0].Score)
	}
}

func TestServiceDoesNotCallAnswererWhenNoHit(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()
	if err := repo.Replace(ctx, []IndexedChunk{
		{
			Document:  Document{ID: "feature:ota:en", Locale: "en", SourceType: "feature", Title: "OTA", URL: "/features/ota", Body: "OTA firmware rollout."},
			ChunkID:   "ota-1",
			Text:      "OTA firmware rollout.",
			Embedding: []float64{1, 0, 0},
		},
	}); err != nil {
		t.Fatalf("replace index: %v", err)
	}
	answerer := &stubAnswerer{}
	service := NewService(repo, &stubEmbedder{vectors: map[string][]float64{"unrelated": {0, 1, 0}}}, answerer, ServiceConfig{MinScore: 0.8})

	result, err := service.Query(ctx, Query{Text: "unrelated", Locale: "en"})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if result.AnswerFound {
		t.Fatalf("answer_found = true, want false: %+v", result)
	}
	if len(result.Sources) != 0 {
		t.Fatalf("sources = %+v, want empty", result.Sources)
	}
	if answerer.calls != 0 {
		t.Fatalf("answerer calls = %d, want 0", answerer.calls)
	}
}

func TestServiceBuildsAnswerFromRetrievedChunksOnly(t *testing.T) {
	repo := newTestRepository(t)
	ctx := context.Background()
	if err := repo.Replace(ctx, []IndexedChunk{
		{
			Document:  Document{ID: "feature:ota:en", Locale: "en", SourceType: "feature", Title: "OTA", URL: "/features/ota", Body: "OTA public body"},
			ChunkID:   "ota-1",
			Text:      "OTA firmware rollout campaign controls.",
			Embedding: []float64{1, 0, 0},
		},
		{
			Document:  Document{ID: "feature:private-cloud:en", Locale: "en", SourceType: "feature", Title: "Private Cloud", URL: "/features/private-cloud", Body: "Private cloud deployment."},
			ChunkID:   "private-1",
			Text:      "Private cloud deployment.",
			Embedding: []float64{0, 1, 0},
		},
	}); err != nil {
		t.Fatalf("replace index: %v", err)
	}
	answerer := &stubAnswerer{}
	service := NewService(repo, &stubEmbedder{vectors: map[string][]float64{"ota": {1, 0, 0}}}, answerer, ServiceConfig{MinScore: 0.8, MaxSources: 1})

	result, err := service.Query(ctx, Query{Text: "ota", Locale: "en"})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if !result.AnswerFound {
		t.Fatalf("answer_found = false, want true: %+v", result)
	}
	if answerer.calls != 1 {
		t.Fatalf("answerer calls = %d, want 1", answerer.calls)
	}
	if !strings.Contains(answerer.context, "OTA firmware rollout campaign controls.") {
		t.Fatalf("answer context missing retrieved chunk: %s", answerer.context)
	}
	if strings.Contains(answerer.context, "Private cloud deployment") {
		t.Fatalf("answer context contains non-retrieved chunk: %s", answerer.context)
	}
}

func TestCollectWebsiteDocumentsIncludesFeatureDocsAndManual(t *testing.T) {
	docs, err := CollectWebsiteDocuments(CollectionConfig{
		RepoRoot:    filepath.Join("..", ".."),
		ContentRoot: filepath.Join("..", "..", "content"),
	})
	if err != nil {
		t.Fatalf("collect documents: %v", err)
	}
	ids := map[string]bool{}
	for _, doc := range docs {
		ids[doc.ID] = true
	}
	for _, id := range []string{
		"feature:ota:en",
		"doc:sdks:en",
		"manual:sdk-samples:en",
		"file:README.md:en",
	} {
		if !ids[id] {
			t.Fatalf("missing collected document %s", id)
		}
	}
}

func newTestRepository(t *testing.T) *Repository {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repo := NewRepository(db)
	if err := repo.Init(context.Background()); err != nil {
		t.Fatalf("init repository: %v", err)
	}
	return repo
}
