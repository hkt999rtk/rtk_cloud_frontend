package search

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo     *Repository
	embedder Embedder
	answerer Answerer
	cfg      ServiceConfig
}

func NewService(repo *Repository, embedder Embedder, answerer Answerer, cfg ServiceConfig) *Service {
	if cfg.MinScore <= 0 {
		cfg.MinScore = 0.35
	}
	if cfg.MaxSources <= 0 {
		cfg.MaxSources = 5
	}
	if strings.TrimSpace(cfg.NoHitText) == "" {
		cfg.NoHitText = "No matching documentation was found."
	}
	return &Service{repo: repo, embedder: embedder, answerer: answerer, cfg: cfg}
}

func (s *Service) Query(ctx context.Context, query Query) (Result, error) {
	text := strings.TrimSpace(query.Text)
	if text == "" {
		return Result{}, fmt.Errorf("query is required")
	}
	if s == nil || s.repo == nil || s.embedder == nil || s.answerer == nil {
		return Result{}, fmt.Errorf("search service is not configured")
	}
	embeddings, err := s.embedder.Embed(ctx, []string{text})
	if err != nil {
		return Result{}, err
	}
	if len(embeddings) != 1 {
		return Result{}, fmt.Errorf("unexpected embedding count")
	}
	results, err := s.repo.Search(ctx, embeddings[0], s.cfg.MaxSources)
	if err != nil {
		return Result{}, err
	}
	sources := filterSources(results, s.cfg.MinScore, s.cfg.MaxSources)
	if len(sources) == 0 {
		return Result{AnswerFound: false, Answer: noHitText(query.Locale), Sources: []Source{}}, nil
	}
	answer, err := s.answerer.Answer(ctx, AnswerRequest{
		Query:   text,
		Locale:  query.Locale,
		Context: buildAnswerContext(sources),
		Sources: sources,
	})
	if err != nil {
		return Result{}, err
	}
	return Result{AnswerFound: true, Answer: strings.TrimSpace(answer), Sources: sources}, nil
}

func filterSources(results []Source, minScore float64, maxSources int) []Source {
	out := make([]Source, 0, maxSources)
	for _, source := range results {
		if source.Score < minScore {
			continue
		}
		source.Snippet = trimSnippet(source.Snippet, 700)
		out = append(out, source)
		if len(out) >= maxSources {
			break
		}
	}
	return out
}

func buildAnswerContext(sources []Source) string {
	var b strings.Builder
	for i, source := range sources {
		b.WriteString(fmt.Sprintf("[Source %d]\nTitle: %s\nURL: %s\nContent: %s\n\n", i+1, source.Title, source.URL, source.Snippet))
	}
	return strings.TrimSpace(b.String())
}

func noHitText(locale string) string {
	switch locale {
	case "zh-TW":
		return "查不到相關文件。"
	case "zh-CN":
		return "查不到相关文档。"
	default:
		return "No matching documentation was found."
	}
}

func trimSnippet(input string, max int) string {
	trimmed := strings.Join(strings.Fields(input), " ")
	if len([]rune(trimmed)) <= max {
		return trimmed
	}
	runes := []rune(trimmed)
	return string(runes[:max]) + "..."
}
