package search

import "context"

type Document struct {
	ID         string
	Locale     string
	SourceType string
	Title      string
	URL        string
	Body       string
}

type IndexedChunk struct {
	Document  Document
	ChunkID   string
	Text      string
	Embedding []float64
}

type Source struct {
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Snippet    string  `json:"snippet"`
	Locale     string  `json:"locale"`
	SourceType string  `json:"source_type"`
	Score      float64 `json:"score"`
}

type Query struct {
	Text   string
	Locale string
}

type Result struct {
	AnswerFound bool     `json:"answer_found"`
	Answer      string   `json:"answer"`
	Sources     []Source `json:"sources"`
}

type Embedder interface {
	Embed(context.Context, []string) ([][]float64, error)
}

type Answerer interface {
	Answer(context.Context, AnswerRequest) (string, error)
}

type AnswerRequest struct {
	Query   string
	Locale  string
	Context string
	Sources []Source
}

type ServiceConfig struct {
	MinScore   float64
	MaxSources int
	NoHitText  string
}
