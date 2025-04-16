package openai

import (
	"context"
	"net/http"
)

const responsesSuffix = "/responses"

type ResponseRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	MaxTokens       int               `json:"max_tokens,omitempty"`
	Temperature     float32           `json:"temperature,omitempty"`
	TopP            float32           `json:"top_p,omitempty"`
	N               int               `json:"n,omitempty"`
	Stream          bool              `json:"stream,omitempty"`
	Stop            []string          `json:"stop,omitempty"`
	PresencePenalty float32           `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32          `json:"frequency_penalty,omitempty"`
	LogitBias       map[string]int    `json:"logit_bias,omitempty"`
	User            string            `json:"user,omitempty"`
	Seed            *int              `json:"seed,omitempty"`
	Store           bool              `json:"store,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type ResponseChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type ResponseResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []ResponseChoice `json:"choices"`
	Usage   Usage            `json:"usage"`

	httpHeader
}

func (c *Client) CreateResponse(
	ctx context.Context,
	request ResponseRequest,
) (response ResponseResponse, err error) {
	if request.Stream {
		err = ErrCompletionStreamNotSupported
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(responsesSuffix, withModel(request.Model)),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
