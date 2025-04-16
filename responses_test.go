package openai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/internal/test/checks"
)

func TestCreateResponse(t *testing.T) {
	client, server, teardown := setupOpenAITestServer()
	defer teardown()
	server.RegisterHandler("/v1/responses", handleResponsesEndpoint)

	resp, err := client.CreateResponse(context.Background(), openai.ResponseRequest{
		Model: openai.GPT4Dot1,
		Input: "Tell me a three sentence bedtime story about a unicorn.",
	})
	checks.NoError(t, err)
	if resp.ID != "resp-123" {
		t.Errorf("Expected ID resp-123, got %s", resp.ID)
	}
	if resp.Choices[0].Text != "Once upon a time, there was a magical unicorn. The unicorn had a rainbow mane and could fly. The end." {
		t.Errorf("Unexpected response text: %s", resp.Choices[0].Text)
	}
}

func handleResponsesEndpoint(w http.ResponseWriter, r *http.Request) {
	var req openai.ResponseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot decode body: %v", err), http.StatusBadRequest)
		return
	}

	response := openai.ResponseResponse{
		ID:      "resp-123",
		Object:  "response",
		Created: 1677825464,
		Model:   req.Model,
		Choices: []openai.ResponseChoice{
			{
				Text:         "Once upon a time, there was a magical unicorn. The unicorn had a rainbow mane and could fly. The end.",
				Index:        0,
				FinishReason: "stop",
			},
		},
		Usage: openai.Usage{
			PromptTokens:     10,
			CompletionTokens: 25,
			TotalTokens:      35,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
