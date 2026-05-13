package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIClient struct {
	APIKey         string
	EmbeddingModel string
	AnswerModel    string
	BaseURL        string
	HTTPClient     *http.Client
}

func (c OpenAIClient) Embed(ctx context.Context, inputs []string) ([][]float64, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}
	model := c.EmbeddingModel
	if model == "" {
		model = "text-embedding-3-small"
	}
	payload := map[string]any{
		"model":           model,
		"input":           inputs,
		"encoding_format": "float",
	}
	var response struct {
		Data []struct {
			Index     int       `json:"index"`
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := c.post(ctx, "/v1/embeddings", payload, &response); err != nil {
		return nil, err
	}
	out := make([][]float64, len(inputs))
	for _, item := range response.Data {
		if item.Index >= 0 && item.Index < len(out) {
			out[item.Index] = item.Embedding
		}
	}
	for i, embedding := range out {
		if len(embedding) == 0 {
			return nil, fmt.Errorf("OpenAI embedding response missing vector for input %d", i)
		}
	}
	return out, nil
}

func (c OpenAIClient) Answer(ctx context.Context, req AnswerRequest) (string, error) {
	if strings.TrimSpace(c.APIKey) == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is required")
	}
	model := c.AnswerModel
	if model == "" {
		model = "gpt-4.1-mini"
	}
	payload := map[string]any{
		"model": model,
		"instructions": strings.Join([]string{
			"You answer questions about Realtek Connect+ documentation.",
			"Use only the provided source context.",
			"If the answer is not in the source context, say that no matching documentation was found.",
			"Keep the answer concise and cite source titles inline when useful.",
		}, "\n"),
		"input":             fmt.Sprintf("Locale: %s\nQuestion: %s\n\nSource context:\n%s", req.Locale, req.Query, req.Context),
		"max_output_tokens": 500,
		"store":             false,
	}
	var response struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := c.post(ctx, "/v1/responses", payload, &response); err != nil {
		return "", err
	}
	if strings.TrimSpace(response.OutputText) != "" {
		return response.OutputText, nil
	}
	for _, item := range response.Output {
		for _, content := range item.Content {
			if strings.TrimSpace(content.Text) != "" {
				return content.Text, nil
			}
		}
	}
	return "", fmt.Errorf("OpenAI response did not include text")
}

func (c OpenAIClient) post(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("OpenAI API %s failed: %s", path, strings.TrimSpace(string(responseBody)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
